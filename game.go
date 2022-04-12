package aegirdungeons

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

const (
	AlgoExplorerTestnet = "https://testnet.algoexplorerapi.io"
	AlgoExplorerMainnet = "https://algoexplorerapi.io"

	NewAlgoExplorerTestnet = "https://node.testnet.algoexplorerapi.io"
	NewAlgoExplorerMainnet = "https://node.algoexplorerapi.io"

	NewAlgoExplorerIndexerTestnet = "https://algoindexer.testnet.algoexplorerapi.io"
	NewAlgoExplorerIndexerMainnet = "https://algoindexer.algoexplorerapi.io"

	IndexerTestnet = "https://testnet.algoexplorerapi.io/idx2"
	IndexerMainnet = "https://algoexplorerapi.io/idx2"

	TestnetDungeonGoldTokenID = 46192450
	MainnetDungeonGoldTokenID = 411521263

	TestnetAegirTokenID = 51294904
	MainnetAegirTokenID = 453816186
)

const (
	DAGRemoji                = "<:dAGR:930481439605149736>"
	GuildID                  = "883754257612943371"
	AdminRoleID              = "903319921679818782"
	DungeoneerRoleID         = "911491170981527602"
	GeneralChatMention       = "<#903284194921295893>"
	AdminRoleMention         = "<@&" + AdminRoleID + ">"
	DungeoneerRoleMention    = "<@&" + DungeoneerRoleID + ">"
	DecimalAdjustment        = 1000000
	DefaultRoundCooldown     = 10 * time.Second
	DefaultPower             = 10
	DefaultAttackCooldown    = 5 * time.Second
	DefaultMinFee            = 0.001 * 1000000
	DefaultRewardAmount      = 10 * DecimalAdjustment
	DefaultMaxEnemyHealth    = 4000
	DefaultMinEnemyHealth    = 1500
	DefaultRewardRatePerHour = 0.000015
	DefaultMaxRewards        = 5000
	DefaultMaxRewards2       = 100
)

const (
	RegisterCommand = "!register"
	AttackCommand   = "!attack"
)

var (
	ErrInvalidCommand = errors.New("game: invalid command")
)

var (
	Limbo      uint32 = 0
	PlayerTurn uint32 = 1
)

// Message ...
type Message struct {
	DiscordUser *discordgo.User
	Content     string
	Roles       []string
}

// Game ...
type Game struct {
	GameConfig

	am                 *AccountManager
	logger             *logrus.Logger
	discordSession     *discordgo.Session
	dungeonGoldTokenID uint64
	lastEnemyTime      time.Time
	messageBuilder     strings.Builder

	bonusTokenName     string
	bonusTokenDecimals float64

	rewardTokenName2     string
	rewardTokenDecimals2 float64

	startRoundChan chan struct{}
	state          uint32

	DungeonInChan chan Message
	currentEnemy  *Enemy
}

// NewGame ...
func NewGame(logger *logrus.Logger, cfg GameConfig) (*Game, error) {
	dg, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		return nil, err
	}

	asaID := uint64(TestnetDungeonGoldTokenID)
	if cfg.MainnetEnabled {
		asaID = MainnetDungeonGoldTokenID
	}

	am, err := NewAccountManager(cfg, asaID)
	if err != nil {
		return nil, err
	}

	g := &Game{}
	g.GameConfig = cfg
	g.dungeonGoldTokenID = asaID
	g.am = am
	g.logger = logger
	g.discordSession = dg
	g.DungeonInChan = make(chan Message)
	g.startRoundChan = make(chan struct{})

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		msg := Message{
			DiscordUser: m.Author,
			Content:     m.Content,
		}
		if m.Member != nil {
			msg.Roles = m.Member.Roles
		}
		if m.ChannelID == g.DungeonDiscordChannelID {
			g.DungeonInChan <- msg
		}
	})
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	if err := dg.Open(); err != nil {
		return nil, err
	}

	return g, nil
}

// Play ...
func (g *Game) Play(ctx context.Context) {
	for {
		select {
		case <-g.startRoundChan:
			g.StartRound(ctx)
		case msg := <-g.DungeonInChan:
			if err := g.HandleDungeonMsg(ctx, msg); err != nil {
				g.logger.Error(err)
			}
		case <-ctx.Done():
			g.Close()
			return
		}
	}
}

// SetBonusAssetDetails ...
func (g *Game) SetBonusAssetDetails(ctx context.Context) error {
	if g.BonusTokenID != 0 {
		asset, err := g.am.FetchAssetInfo(ctx, g.BonusTokenID)
		if err != nil {
			return err
		}
		g.bonusTokenDecimals = math.Pow10(int(asset.Params.Decimals))
		g.bonusTokenName = asset.Params.UnitName
	}

	return nil
}

// SetReward2AssetDetails ...
func (g *Game) SetReward2AssetDetails(ctx context.Context) error {
	if g.RewardAssetID2 != 0 {
		asset, err := g.am.FetchAssetInfo(ctx, g.RewardAssetID2)
		if err != nil {
			return err
		}
		g.rewardTokenDecimals2 = math.Pow10(int(asset.Params.Decimals))
		g.rewardTokenName2 = asset.Params.UnitName
	}

	return nil
}

// Send ...
func (g *Game) Send(msg string) error {
	ms := discordgo.MessageSend{
		Content: msg,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{
				discordgo.AllowedMentionTypeRoles,
				discordgo.AllowedMentionTypeUsers,
				discordgo.AllowedMentionTypeEveryone,
			},
		},
	}

	_, err := g.discordSession.ChannelMessageSendComplex(g.DungeonDiscordChannelID, &ms)
	return err
}

// StartRoundCoolDown ...
func (g *Game) StartRoundCoolDown() {
	g.state = Limbo

	m := fmt.Sprintf("---\n"+
		"⏳ Waiting For next Round to start in %v ⏳\n"+
		"---", g.RoundCoolDown)
	if err := g.Send(m); err != nil {
		g.logger.Errorf("start_round_cool_down: %v\n", err)
	}
	g.lastEnemyTime = time.Now()

	time.Sleep(g.RoundCoolDown - (time.Minute * 2))
	if err := g.InitializeAndWait(context.Background(), (time.Minute * 2)); err != nil {
		g.logger.Errorf("start_round: initialize_game:", err)
	}
	g.startRoundChan <- struct{}{}
}

// InitializeAndWait ...
func (g *Game) InitializeAndWait(ctx context.Context, d time.Duration) error {
	// Game Starting Warning
	m := fmt.Sprintf("%s %s Warning\n", DungeoneerRoleMention, d)
	if err := g.Send(m); err != nil {
		g.logger.Errorf("start_round_cooldown: %v\n", err)
	}

	g.am.ClearCaches()
	if err := g.am.InitializeCaches(ctx); err != nil {
		return err
	}
	time.Sleep(d)

	return nil
}

// StartRound ...
func (g *Game) StartRound(ctx context.Context) {
	if err := g.SetBonusAssetDetails(ctx); err != nil {
		g.logger.Fatal(err)
	}
	if err := g.SetReward2AssetDetails(ctx); err != nil {
		g.logger.Fatal(err)
	}

	g.CreateEnemy(ctx)

	m := fmt.Sprintf("%s 📯 STARTING ROUND 📯\n", DungeoneerRoleMention)
	m += g.currentEnemy.String()
	g.state = PlayerTurn
	if err := g.Send(m); err != nil {
		g.logger.Errorf("start_round: %v\n", err)
	}

}

// CreateEnemy ...
func (g *Game) CreateEnemy(ctx context.Context) {
	e := RandomEnemy(g.MinEnemyHealth, g.MaxEnemyHealth)
	if g.RewardAmount != 0 {
		e.RewardAmount = g.RewardAmount
	} else {
		rewardAmount, err := g.am.RewardAmount(ctx, g.RewardRate)
		if err != nil {
			g.logger.Errorf("create_enemy:", err)
			rewardAmount = DefaultRewardAmount
		}
		e.RewardAmount = rewardAmount
	}

	remaining, err := g.am.AssetReserves(ctx, g.BonusTokenID)
	if err != nil {
		g.logger.Errorf("create_enemy:", err)
	}
	if err == nil && g.BonusRewards > 0 && remaining > 0 {
		e.BonusRewardAmount = g.BonusRewards * g.bonusTokenDecimals

		accs, err := g.am.ValidAsaAccounts(ctx, g.BonusTokenID)
		if err != nil {
			g.logger.Errorf("create_enemy:", err)
		}
		for _, acc := range accs {
			e.BonusLooters[acc] = 0
		}
	}

	remaining, err = g.am.AssetReserves(ctx, g.RewardAssetID2)
	if err != nil {
		g.logger.Errorf("create_enemy:", err)
	}
	if err == nil && remaining > 0 {
		e.RewardAmount2 = g.RewardRate2 * g.rewardTokenDecimals2
		if float64(remaining) < g.RewardRate2 {
			e.RewardAmount2 = float64(remaining)
		}

		accs, err := g.am.ValidAsaAccounts(ctx, g.RewardAssetID2)
		if err != nil {
			g.logger.Errorf("create_enemy:", err)
		}
		for _, acc := range accs {
			e.AddLooter2(acc)
		}
	}

	g.currentEnemy = &e
}

// HasAdminRole ...
func HasAdminRole(roles []string) bool {
	for i := 0; i < len(roles); i++ {
		if roles[i] == AdminRoleID {
			return true
		}
	}
	return false
}

// HandleDungeonMsg ...
func (g *Game) HandleDungeonMsg(ctx context.Context, msg Message) error {
	command := strings.ToLower(strings.TrimSpace(msg.Content))
	if strings.HasPrefix(command, RegisterCommand) {
		if err := g.Register(ctx, msg); err != nil {
			return err
		}
		return nil
	}

	if err := g.am.ValidateAccount(ctx, msg.DiscordUser.ID); err != nil {
		return g.HandleAccountErr(msg, err)
	}

	switch g.state {
	case Limbo:
		if command == AttackCommand {
			dt := g.lastEnemyTime.Add(g.RoundCoolDown).Sub(time.Now())
			if dt < 0 {
				return nil
			}
			m := fmt.Sprintf("%s next round will start in %s.\n", msg.DiscordUser.Mention(), dt.Truncate(time.Second).String())
			if err := g.Send(m); err != nil {
				g.logger.Error("attack_enemy:", err)
			}
			return nil
		}
	case PlayerTurn:
		if command == AttackCommand {
			return g.AttackEnemy(ctx, msg)
		}
	default:
		return fmt.Errorf("handle_msg: unknown game state %q\n", g.state)
	}

	return nil
}

// AttackEnemy ...
func (g *Game) AttackEnemy(ctx context.Context, msg Message) error {
	address, err := g.am.Address(ctx, msg.DiscordUser.ID)
	if err != nil {
		return err
	}

	legends, err := g.am.Legends(ctx, address)
	if err != nil {
		return err
	}
	emojis := g.am.LegendEmojis(legends)
	power := g.am.GetPower(legends)

	if g.am.IsDuplicateUser(ctx, msg.DiscordUser.ID, address) {
		m := fmt.Sprintf("%s looks like multiple users are using this account!\n Please speak with an %s\n", msg.DiscordUser.Mention(), AdminRoleMention)
		return g.Send(m)
	}
	g.am.usersToAccount[msg.DiscordUser.ID] = address
	g.am.LastAttack[address] = time.Now()

	if g.currentEnemy == nil {
		return nil
	}
	g.currentEnemy.Damage(address, float64(power))

	m := fmt.Sprintf("%s Dealt %v damage! %s ⚔ %s, %v HP Remaining!\n", msg.DiscordUser.Mention(), power, strings.Join(emojis, " "), g.currentEnemy.Name, g.currentEnemy.Health)
	if g.currentEnemy.Health <= 0 {
		g.currentEnemy.Damage(address, float64(power))
		m = fmt.Sprintf("%s 🎉🎉🎉 %s Dealt the final hit! X2 DAMAGE!! 🎉🎉🎉\n", msg.DiscordUser.Mention(), strings.Join(emojis, " "))
	}
	if g.currentEnemy.isFirstHit {
		g.currentEnemy.isFirstHit = false
		g.currentEnemy.Damage(address, float64(power))
		m = fmt.Sprintf("%s 🎉🎉🎉 %s Dealt the FIRST hit! X2 DAMAGE!! 🎉🎉🎉\n", msg.DiscordUser.Mention(), strings.Join(emojis, " "))
	}
	if _, err := g.messageBuilder.WriteString(m); err != nil {
		g.logger.Error("attack_enemy:", err)
	}

	if g.currentEnemy.Health <= 0 || g.messageBuilder.Len() > 1600 {
		if err := g.Send(g.messageBuilder.String()); err != nil {
			g.logger.Error("attack_enemy:", err)
		}
		g.messageBuilder.Reset()
	}

	if g.currentEnemy.Health <= 0 {
		rewardLedger := g.currentEnemy.DropLedger()
		rewardAmount := g.currentEnemy.RewardAmount / DecimalAdjustment
		rewardLedger2 := g.currentEnemy.DropLedger2()
		rewardAmount2 := g.currentEnemy.RewardAmount2 / g.rewardTokenDecimals2
		bonusLedger := g.currentEnemy.DropBonusLedger()

		g.currentEnemy = nil
		defer func() {
			go g.StartRoundCoolDown()
		}()
		if _, err := g.am.SendLoot(ctx, g.dungeonGoldTokenID, rewardLedger); err != nil {
			return err
		}
		m := "☠ Enemy Killed ☠ \n"
		m += fmt.Sprintf("%s: %vdAGR\n", DAGRemoji, rewardAmount)
		if len(bonusLedger) > 0 {
			if _, err := g.am.SendLoot(ctx, g.BonusTokenID, bonusLedger); err != nil {
				return err
			}
			m += fmt.Sprintf("💰: %v%s\n", len(bonusLedger)*int(g.BonusRewards), g.bonusTokenName)
		}
		if len(rewardLedger2) > 0 {
			if _, err := g.am.SendLoot(ctx, g.RewardAssetID2, rewardLedger2); err != nil {
				return err
			}
			m += fmt.Sprintf("💰: %v%s\n", rewardAmount2, g.rewardTokenName2)
		}
		m += fmt.Sprintf("👥: %v\n", len(rewardLedger))
		if err := g.Send(m); err != nil {
			return err
		}
	}

	return nil
}

// Register ...
func (g *Game) Register(ctx context.Context, msg Message) error {
	address := strings.TrimPrefix(msg.Content, RegisterCommand)
	address = strings.TrimSpace(address)
	if err := g.am.Register(ctx, msg.DiscordUser.ID, address); err != nil {
		if err == ErrMultipleRegisteredUsers {
			m := fmt.Sprintf("%s looks like multiple users are using this account!\n Please speak with an %s\n", msg.DiscordUser.Mention(), AdminRoleMention)
			return g.Send(m)
		}

		g.logger.Warnf("register: %v, %q\n", err, address)
		return g.Send("Invalid Address. Usage `!register <58 length address>`")
	}
	if err := g.discordSession.GuildMemberRoleAdd(GuildID, msg.DiscordUser.ID, DungeoneerRoleID); err != nil {
		g.logger.Errorf("register: %v\n", err)
	}
	m := fmt.Sprintf("%s Your account has been registered.", msg.DiscordUser.Mention())
	return g.Send(m)
}

// HandleAccountErr
func (g *Game) HandleAccountErr(msg Message, inErr error) error {
	m := fmt.Sprintf("%s your account is not able to play.\n", msg.DiscordUser.Mention())
	switch inErr {
	case ErrUserNotRegistered:
		m += "Ensure your account is registered. Use `!register <58 length address>`"
	case ErrNoLegends:
		m += "Ensure your account holds a Legend. https://www.nftexplorer.app/sellers?creator=TDRZOHTX4LFSYX3LUXL5WME47UUOUSAYHK4DL23B6R4OLAYJTYBJAGGNJI"
	case ErrNotOptedIn:
		m += fmt.Sprintf("Ensure your account has opted into ASA#%v; Press the Accept ASA button https://www.randgallery.com/algo-collection/?address=%v", g.dungeonGoldTokenID, g.dungeonGoldTokenID)
	default:
		m += "The reason will have to be investigated"
		g.logger.Error("handle_account:", inErr)
	}

	return g.Send(m)
}

// Close ...
func (g *Game) Close() {
	if err := g.Send("🛑 Aegir Dungeons is temporarily offline 🛑"); err != nil {
		g.logger.Error("closing game:", err)
	}
	if err := g.am.Close(); err != nil {
		g.logger.Error("closing account_manager:", err)
	}
	close(g.startRoundChan)
	close(g.DungeonInChan)
	if err := g.discordSession.Close(); err != nil {
		g.logger.Error("closing discord:", err)
	}
	g.logger.Info("stopping game")
}
