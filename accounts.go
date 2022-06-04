package aegirdungeons

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
)

const (
	ValidAddressLength = 58
)

var (
	ErrNoUserFound             = errors.New("accounts: no user found")
	ErrNoLegends               = errors.New("accounts: no legends in account")
	ErrNotOptedIn              = errors.New("accounts: account has not opted into rewards")
	ErrInvalidAddress          = errors.New("accounts: invalid address")
	ErrAssetNotFound           = errors.New("accounts: asset not found")
	ErrMultipleRegisteredUsers = errors.New("accounts: multiple users found with same address")
	ErrDiscordUserIDNotFound   = errors.New("accounts: discord user id not found")
)

// Account ...
type Account struct {
	DiscordID string
	Address   string
}

// AccountStore ...
type AccountStore interface {
	Set(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) (string, error)
	Close() error
}

// AccountManager ...
type AccountManager struct {
	ac *algod.Client
	ai *indexer.Client

	accountStore       AccountStore
	LastAttack         map[string]time.Time
	AsaToClass         map[uint64]uint64
	AsaToEmoji         map[uint64]string
	dungeonGoldTokenID uint64

	legendCache     map[string][]uint64
	usersToAccount  map[string]string
	accountToUserID map[string]string
	optInCache      map[string]struct{}
	hasLegendCache  map[string]struct{}
	miniLevelCache  map[uint64]uint64

	publicKey string
	mnemonic  string
}

// NewAccountManager ...
func NewAccountManager(cfg GameConfig, asaID uint64) (*AccountManager, error) {
	asaToClass := TestnetAsaToClass
	asaToEmoji := TestnetAsaToEmoji
	address := NewAlgoExplorerTestnet
	indexerAddress := NewAlgoExplorerIndexerTestnet
	if cfg.MainnetEnabled {
		asaToClass = MainnetAsaToClass
		asaToEmoji = MainnetAsaToEmoji
		address = NewAlgoExplorerMainnet
		indexerAddress = NewAlgoExplorerIndexerMainnet
	}
	ac, err := algod.MakeClient(address, "")
	if err != nil {
		return nil, fmt.Errorf("new_account_manager: algod: make_client: %v", err)
	}

	ai, err := indexer.MakeClient(indexerAddress, "")
	if err != nil {
		return nil, fmt.Errorf("new_account_manager: indexer: make_client: %v", err)
	}

	var as AccountStore
	if cfg.RedisAddress != "" {
		as, err = NewRedisStore(cfg.RedisAddress)
		if err != nil {
			return nil, fmt.Errorf("new_account_manager: new_redis_store: %v", err)
		}
	} else {
		as = NewInMemoryStore()
	}

	return &AccountManager{
		dungeonGoldTokenID: asaID,
		ac:                 ac,
		ai:                 ai,
		accountStore:       as,
		LastAttack:         map[string]time.Time{},
		AsaToClass:         asaToClass,
		AsaToEmoji:         asaToEmoji,
		usersToAccount:     map[string]string{},
		accountToUserID:    map[string]string{},
		publicKey:          cfg.BankPublic,
		mnemonic:           cfg.BankMnemonic,

		legendCache:    make(map[string][]uint64),
		optInCache:     make(map[string]struct{}),
		hasLegendCache: make(map[string]struct{}),
		miniLevelCache: make(map[uint64]uint64),
	}, nil
}

// InitializeCaches ...
func (am *AccountManager) InitializeCaches(ctx context.Context) error {
	if err := am.InitializeLegendsCache(ctx); err != nil {
		return fmt.Errorf("initialize_caches: %v", err)
	}

	if err := am.InitializeOptInCache(ctx, am.dungeonGoldTokenID); err != nil {
		return fmt.Errorf("initialize_caches: %v", err)
	}

	return nil
}

// RewardAmount ...
func (am *AccountManager) RewardAmount(ctx context.Context, rewardRate float64) (float64, error) {
	_, acc, err := am.ai.LookupAccountByID(am.publicKey).Do(ctx)
	// acc, err := am.ac.AccountInformation(am.publicKey).Do(ctx)
	if err != nil {
		return 0, err
	}

	for _, a := range acc.Assets {
		if a.AssetId == am.dungeonGoldTokenID {
			amt := float64(a.Amount) * rewardRate
			if amt > (4000 * DecimalAdjustment) {
				amt = DefaultRewardAmount
			}

			return amt, nil
		}
	}

	return 0, ErrNotOptedIn
}

// LegendEmojis ...
func (am *AccountManager) LegendEmojis(asaIDs []uint64) []string {
	var emojis []string
	for _, asaID := range asaIDs {
		emoji := am.AsaToEmoji[asaID]
		emojis = append(emojis, emoji)
	}

	return emojis
}

// Address ...
func (am *AccountManager) Address(ctx context.Context, userID string) (string, error) {
	address, ok := am.usersToAccount[userID]
	if ok {
		return address, nil
	}

	address, err := am.accountStore.Get(ctx, userID)
	if err != nil {
		return "", err
	}
	am.usersToAccount[userID] = address

	return address, nil
}

// IsDuplicateUser ...
func (am *AccountManager) IsDuplicateUser(ctx context.Context, userID, address string) bool {
	registeredUserID, ok := am.accountToUserID[address]
	if ok && registeredUserID != userID {
		return true
	}
	registeredUserID, err := am.accountStore.Get(ctx, fmt.Sprintf("WAL-", address))
	if err != nil || registeredUserID == "" {
		am.accountStore.Set(ctx, fmt.Sprintf("WAL-", address), userID)
		return false
	}
	if registeredUserID != userID {
		return true
	}

	return false
}

// Register ...
func (am *AccountManager) Register(ctx context.Context, userID, address string) error {
	// Check Address Length
	if len(address) != ValidAddressLength {
		return ErrInvalidAddress
	}
	if am.IsDuplicateUser(ctx, userID, address) {
		return ErrMultipleRegisteredUsers
	}

	if err := am.accountStore.Set(ctx, fmt.Sprintf("WAL-", address), userID); err != nil {
		return err
	}

	return am.accountStore.Set(ctx, userID, address)
}

// ValidateAccount ...
func (am *AccountManager) ValidateAccount(ctx context.Context, userID string) error {
	// Check if Account Registered
	address, err := am.Address(ctx, userID)
	if err != nil {
		return err
	}

	// Check if Opted into dugeon rewards
	if err := am.ValidateOptIn(ctx, address, am.dungeonGoldTokenID); err != nil {
		return err
	}

	// Check if Holds legends
	if err := am.ValidateLegends(ctx, address); err != nil {
		return err
	}

	return nil
}

// ValidateOptIn ...
func (am *AccountManager) ValidateOptIn(ctx context.Context, address string, asaID uint64) error {
	key := strconv.FormatUint(asaID, 10) + "/" + address
	if _, ok := am.optInCache[key]; ok {
		return nil
	}

	// acc, err := am.ac.AccountInformation(address).Do(ctx)
	_, acc, err := am.ai.LookupAccountByID(address).Do(ctx)
	if err != nil {
		return err
	}

	for _, a := range acc.Assets {
		if a.AssetId == asaID {
			am.optInCache[key] = struct{}{}
			return nil
		}
	}

	return ErrNotOptedIn
}

// InitializeOptInCache ...
func (am *AccountManager) InitializeOptInCache(ctx context.Context, asaID uint64) error {
	nextToken := ""

	for {
		resp, err := am.ai.LookupAssetBalances(asaID).CurrencyGreaterThan(0).Limit(1000).NextToken(nextToken).Do(ctx)
		if err != nil {
			return err
		}

		for _, bal := range resp.Balances {
			key := strconv.FormatUint(asaID, 10) + "/" + bal.Address
			am.optInCache[key] = struct{}{}
		}

		nextToken = resp.NextToken
		if nextToken == "" {
			break
		}
	}

	return nil
}

// InitializeLegendsCache ...
func (am *AccountManager) InitializeLegendsCache(ctx context.Context) error {
	addressToMinis := make(map[string]uint64)

	for asaID, class := range am.AsaToClass {
		nextToken := ""

		for {
			resp, err := am.ai.LookupAssetBalances(asaID).CurrencyGreaterThan(0).Limit(1000).NextToken(nextToken).Do(ctx)
			if err != nil {
				return err
			}

			for _, bal := range resp.Balances {
				if bal.Amount <= 0 {
					continue
				}

				am.hasLegendCache[bal.Address] = struct{}{}
				if class != Mini {
					am.legendCache[bal.Address] = append(am.legendCache[bal.Address], asaID)
					continue
				}

				arc69, err := am.FetchArc69Metadata(ctx, asaID)
				if err != nil {
					return err
				}

				miniLevel, err := strconv.ParseUint(arc69.Properties.Level, 10, 64)
				if err != nil {
					return err
				}

				am.miniLevelCache[asaID] = miniLevel

				currentMaxMini := addressToMinis[bal.Address]
				currentMaxMiniLevel := am.miniLevelCache[currentMaxMini]
				if miniLevel >= currentMaxMiniLevel {
					addressToMinis[bal.Address] = asaID
				}

			}

			nextToken = resp.NextToken
			if nextToken == "" {
				break
			}
		}
	}

	for address, miniAsaID := range addressToMinis {
		am.legendCache[address] = append(am.legendCache[address], miniAsaID)
	}

	return nil
}

// ValidateLegends ...
func (am *AccountManager) ValidateLegends(ctx context.Context, address string) error {
	if _, ok := am.hasLegendCache[address]; ok {
		return nil
	}

	// acc, err := am.ac.AccountInformation(address).Do(ctx)
	_, acc, err := am.ai.LookupAccountByID(address).Do(ctx)
	if err != nil {
		return err
	}

	for _, a := range acc.Assets {
		if a.Amount == 0 {
			continue
		}
		if _, ok := am.AsaToClass[a.AssetId]; ok {
			am.hasLegendCache[address] = struct{}{}
			return nil
		}
	}

	return ErrNoLegends
}

// ValidateGovernance ...
func (am *AccountManager) ValidateGovernance(ctx context.Context, userID string) error {
	// Check if Account Registered
	address, err := am.Address(ctx, userID)
	if err != nil {
		return err
	}

	// Check if Opted into dugeon rewards
	if err := am.ValidateOptIn(ctx, address, am.dungeonGoldTokenID); err != nil {
		return err
	}

	return nil
}

// AssetReserves ...
func (am *AccountManager) AssetReserves(ctx context.Context, asaID uint64) (uint64, error) {
	// acc, err := am.ac.AccountInformation(am.publicKey).Do(ctx)
	_, acc, err := am.ai.LookupAccountByID(am.publicKey).Do(ctx)
	if err != nil {
		return 0, err
	}

	for _, a := range acc.Assets {
		if a.AssetId == asaID {
			return a.Amount, nil
		}
	}

	return 0, ErrAssetNotFound
}

// ClearCaches ...
func (am *AccountManager) ClearCaches() {
	am.legendCache = make(map[string][]uint64)
	am.usersToAccount = make(map[string]string)
	am.hasLegendCache = make(map[string]struct{})
	am.optInCache = make(map[string]struct{})
	am.miniLevelCache = make(map[uint64]uint64)
}

// Legends ...
func (am *AccountManager) Legends(ctx context.Context, address string) ([]uint64, error) {
	if v, ok := am.legendCache[address]; ok {
		return v, nil
	}

	// acc, err := am.ac.AccountInformation(address).Do(ctx)
	_, acc, err := am.ai.LookupAccountByID(address).Do(ctx)
	if err != nil {
		return nil, err
	}

	var legends []uint64
	for _, a := range acc.Assets {
		if a.Amount == 0 {
			continue
		}
		if _, ok := am.AsaToClass[a.AssetId]; ok {
			legends = append(legends, a.AssetId)
		}
	}
	am.legendCache[address] = legends

	return legends, nil
}

// LegendHolderDiscordIDs ...
func (am *AccountManager) LegendHolderDiscordIDs(ctx context.Context) ([]string, error) {
	res := make([]string, 0, len(am.legendCache))
	for address := range am.legendCache {
		userID, err := am.GetDiscordUserID(ctx, address)
		if err != nil {
			if err == ErrDiscordUserIDNotFound {
				continue
			}
			return nil, err
		}
		res = append(res, userID)
	}

	return res, nil
}

// GetDiscordUserID ...
func (am *AccountManager) GetDiscordUserID(ctx context.Context, address string) (string, error) {
	userID, ok := am.accountToUserID[address]
	if ok {
		return userID, nil
	}

	userID, err := am.accountStore.Get(ctx, fmt.Sprintf("WAL-", address))
	if err != nil {
		return "", err
	}
	if userID == "" {
		return "", ErrDiscordUserIDNotFound
	}
	am.accountToUserID[address] = userID

	return userID, nil
}

// TimeSinceLastAttack ...
func (am *AccountManager) TimeSinceLastAttack(address string) time.Duration {
	return time.Now().Sub(am.LastAttack[address])
}

// SendLoot ...
func (am *AccountManager) SendLoot(ctx context.Context, asaID uint64, ledger map[string]float64) ([]string, error) {
	txParams, err := am.ac.SuggestedParams().Do(ctx)
	if err != nil {
		return nil, err
	}

	sk, err := mnemonic.ToPrivateKey(am.mnemonic)
	if err != nil {
		return nil, err
	}

	genHashstr := base64.StdEncoding.EncodeToString(txParams.GenesisHash)
	firstValidRound := txParams.FirstRoundValid
	lastValidRound := txParams.LastRoundValid
	genesisID := txParams.GenesisID
	txs := []types.Transaction{}

	for destAddr, amount := range ledger {
		tx, err := transaction.MakeAssetTransferTxnWithFlatFee(
			am.publicKey,
			destAddr, "", uint64(amount),
			DefaultMinFee,
			uint64(firstValidRound),
			uint64(lastValidRound),
			nil,
			genesisID,
			genHashstr,
			asaID)
		if err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}

	var ptxs []string
	var cur, end int
	for end != len(txs) {
		signedGroup := []byte{}
		end += 15
		if end > len(txs) {
			end = len(txs)
		}
		tempTxs := txs[cur:end]
		gid, err := crypto.ComputeGroupID(tempTxs)
		if err != nil {
			return ptxs, err
		}

		for _, tx := range tempTxs {
			tx.Group = gid
			_, stx, err := crypto.SignTransaction(sk, tx)
			if err != nil {
				return ptxs, err
			}
			signedGroup = append(signedGroup, stx...)
		}

		ptx, err := am.ac.SendRawTransaction(signedGroup).Do(ctx)
		if err != nil {
			return ptxs, fmt.Errorf("send_loot:", err)
		}
		ptxs = append(ptxs, ptx)
		cur += 15
	}

	return ptxs, nil
}

// GetPower ...
func (am *AccountManager) GetPower(legends []uint64) int {
	var power int
	var minis int

	var maxMiniBonus uint64
	for i := 0; i < len(legends); i++ {
		class := am.AsaToClass[legends[i]]
		if class == Mini {
			minis = 1
			minisBonus := am.miniLevelCache[legends[i]]
			if minisBonus > maxMiniBonus {
				maxMiniBonus = minisBonus
			}
			continue
		}

		power += DefaultPower
	}

	power = power + minis + int(float64(power)*(float64(maxMiniBonus)/100.0))
	return power
}

// ValidAsaAccounts ...
func (am *AccountManager) ValidAsaAccounts(ctx context.Context, asaID uint64) ([]string, error) {
	var res []string
	var nextToken string
	for {
		resp, err := am.ai.LookupAssetBalances(asaID).Limit(1000).NextToken(nextToken).Do(ctx)
		if err != nil {
			return nil, err
		}
		for _, bal := range resp.Balances {
			if bal.Address == am.publicKey {
				// TODO: THIS NEEDS TO CHECK FOR HOW MUCH BONUS IS AVAILABLE
				continue
			}

			res = append(res, bal.Address)
		}
		nextToken = resp.NextToken
		if nextToken == "" {
			break
		}
	}

	return res, nil
}

// FetchArc69Metadata ...
func (am *AccountManager) FetchArc69Metadata(ctx context.Context, asaID uint64) (Arc69, error) {
	_, arc69, err := FetchArc69Metadata(ctx, am.ai, asaID, am.publicKey)
	return arc69, err
}

// FetchAssetInfo ...
func (am *AccountManager) FetchAssetInfo(ctx context.Context, asaID uint64) (models.Asset, error) {
	_, asset, err := am.ai.LookupAssetByID(asaID).Do(ctx)
	return asset, err
}

func (am *AccountManager) Close() error {
	return am.accountStore.Close()
}
