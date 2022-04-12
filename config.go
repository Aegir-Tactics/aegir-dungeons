package aegirdungeons

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// GameConfig ...
type GameConfig struct {
	Port                    string
	BankPublic              string
	BankMnemonic            string
	DiscordToken            string
	DungeonDiscordChannelID string
	RedisAddress            string
	MainnetEnabled          bool
	RoundCoolDown           time.Duration
	AttackCoolDown          time.Duration
	MinEnemyHealth          int
	MaxEnemyHealth          int

	RewardAmount float64
	RewardRate   float64

	RewardAssetID2 uint64
	RewardRate2    float64

	BonusTokenID uint64
	BonusRewards float64
}

// NewGameConfig
func NewGameConfig() (GameConfig, error) {
	gc := GameConfig{}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	gc.Port = port

	bankPublic := os.Getenv("BANK_PUBLIC")
	if bankPublic == "" {
		return gc, fmt.Errorf("$BANK_PUBLIC must be set")
	}
	gc.BankPublic = bankPublic

	bankMnemonic := os.Getenv("BANK_MNEMONIC")
	if bankMnemonic == "" {
		return gc, fmt.Errorf("$BANK_MNEMONIC must be set")
	}
	gc.BankMnemonic = bankMnemonic

	discordToken := os.Getenv("DISCORD_TOKEN")
	if discordToken == "" {
		return gc, fmt.Errorf("$DISCORD_TOKEN must be set")
	}
	gc.DiscordToken = discordToken

	dungeonDiscordChannelID := os.Getenv("DISCORD_CHANNEL_ID")
	if dungeonDiscordChannelID == "" {
		return gc, fmt.Errorf("$DISCORD_CHANNEL_ID must be set")
	}
	gc.DungeonDiscordChannelID = dungeonDiscordChannelID

	roundCoolDown := os.Getenv("ROUND_COOLDOWN")
	gc.RoundCoolDown = DefaultRoundCooldown
	if roundCoolDown != "" {
		d, err := time.ParseDuration(roundCoolDown)
		if err != nil {
			return gc, fmt.Errorf("$ROUND_COOLDOWN error: %v", err)
		}
		gc.RoundCoolDown = d
	}
	gc.RewardRate = DefaultRewardRatePerHour * gc.RoundCoolDown.Hours()

	// TODO REMOVE ATTACK COOLDOWN
	attackCoolDown := os.Getenv("ATTACK_COOLDOWN")
	gc.AttackCoolDown = DefaultAttackCooldown
	if attackCoolDown != "" {
		d, err := time.ParseDuration(attackCoolDown)
		if err != nil {
			return gc, fmt.Errorf("$ATTACK_COOLDOWN error: %v", err)
		}
		gc.AttackCoolDown = d
	}

	deployment := strings.ToLower(os.Getenv("DEPLOYMENT_ENV"))
	if deployment == "prod" || deployment == "production" {
		gc.MainnetEnabled = true
	}

	gc.RedisAddress = os.Getenv("REDIS_URL")

	rewardAmount := os.Getenv("REWARD_AMOUNT")
	if rewardAmount != "" {
		d, err := strconv.ParseFloat(rewardAmount, 64)
		if err != nil {
			return gc, fmt.Errorf("$REWARD_AMOUNT error: %v", err)
		}
		gc.RewardAmount = d * DecimalAdjustment
	}

	minHP := os.Getenv("MIN_ENEMY_HEALTH")
	gc.MinEnemyHealth = DefaultMinEnemyHealth
	if minHP != "" {
		d, err := strconv.Atoi(minHP)
		if err != nil {
			return gc, fmt.Errorf("$MIN_ENEMY_HEALTH error: %v", err)
		}
		gc.MinEnemyHealth = d
	}

	maxHP := os.Getenv("MAX_ENEMY_HEALTH")
	gc.MaxEnemyHealth = DefaultMaxEnemyHealth
	if maxHP != "" {
		d, err := strconv.Atoi(maxHP)
		if err != nil {
			return gc, fmt.Errorf("$MAX_ENEMY_HEALTH error: %v", err)
		}
		gc.MaxEnemyHealth = d
	}

	bonusTokenID := os.Getenv("BONUS_TOKEN_ID")
	if bonusTokenID != "" {
		u, err := strconv.ParseUint(bonusTokenID, 10, 64)
		if err != nil {
			return gc, fmt.Errorf("$BONUS_TOKEN_ID error: %v", err)
		}
		gc.BonusTokenID = u
	}

	bonusRewards := os.Getenv("BONUS_REWARDS")
	if bonusRewards != "" {
		d, err := strconv.ParseFloat(bonusRewards, 64)
		if err != nil {
			return gc, fmt.Errorf("$BONUS_REWARDS error: %v", err)
		}
		gc.BonusRewards = d
	}

	rewardAssetID2 := os.Getenv("REWARD_ASSET_ID_2")
	if rewardAssetID2 != "" {
		u, err := strconv.ParseUint(rewardAssetID2, 10, 64)
		if err != nil {
			return gc, fmt.Errorf("$REWARD_ASSET_ID_2 error: %v", err)
		}
		gc.RewardAssetID2 = u
	}

	rewardRate2 := os.Getenv("REWARD_RATE_2")
	if rewardRate2 != "" {
		d, err := strconv.ParseFloat(rewardRate2, 64)
		if err != nil {
			return gc, fmt.Errorf("$REWARD_RATE_2 error: %v", err)
		}
		gc.RewardRate2 = d
	}

	return gc, nil
}
