package aegirdungeons

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var Enemies = map[int]string{
	0:  "ðđ",
	1:  "ðš",
	2:  "ðŧ",
	3:  "ð",
	4:  "ðŋ",
	5:  "ðą",
	6:  "ðĪĒ",
	7:  "ðĪ",
	8:  "ðū",
	9:  "ðĪ",
	10: "ðĪŠ",
	11: "ð",
	12: "ð―",
	13: "ðĪĄ",
	14: "ð§",
	15: "ð§",
	16: "ð§",
	17: "ð§",
	18: "ð§",
	19: "ð§",
	20: "ðž",
}

// RandomEnermy ...
func RandomEnemy(minHP, maxHP int) Enemy {
	h := rand.Intn(maxHP-minHP) + minHP
	v := rand.Intn(len(Enemies))

	return NewEnemy(Enemies[v], float64(h))
}

// Enemy ...
type Enemy struct {
	Name              string
	Health            float64
	RewardAmount      float64
	RewardAmount2     float64
	BonusRewardAmount float64

	totalDamageDone  float64
	totalDamageDone2 float64

	looters      map[string]float64
	looters2     map[string]float64
	BonusLooters map[string]float64

	isFirstHit bool
}

// NewEnemy ...
func NewEnemy(name string, health float64) Enemy {
	return Enemy{
		Name:         name,
		Health:       health,
		BonusLooters: map[string]float64{},
		looters:      map[string]float64{},
		looters2:     map[string]float64{},
		isFirstHit:   true,
	}
}

// AddLooter2 ...
func (e *Enemy) AddLooter2(address string) {
	e.looters2[address] = 0
}

// Damage ...
func (e *Enemy) Damage(address string, amount float64) {
	e.Health -= amount
	e.looters[address] += amount
	e.totalDamageDone += amount

	if _, ok := e.looters2[address]; ok {
		e.looters2[address] += amount
		e.totalDamageDone2 += amount
	}
}

// String ...
func (e *Enemy) String() string {
	return fmt.Sprintf(`A %s has appeared! %v HP`, e.Name, e.Health)
}

// DropLedger ...
func (e *Enemy) DropLedger() map[string]float64 {
	res := map[string]float64{}

	for looter, damage := range e.looters {
		amt := (damage / e.totalDamageDone) * e.RewardAmount
		if amt > 0 {
			res[looter] = amt
		}
	}

	return res
}

// DropLedger2...
func (e *Enemy) DropLedger2() map[string]float64 {
	res := map[string]float64{}

	for looter, damage := range e.looters2 {
		amt := (damage / e.totalDamageDone2) * e.RewardAmount2
		if amt > 0 {
			res[looter] = amt
		}
	}

	return res
}

// DropBonusLedger ...
func (e *Enemy) DropBonusLedger() map[string]float64 {
	res := map[string]float64{}

	for looter, damage := range e.looters {
		if _, ok := e.BonusLooters[looter]; !ok {
			continue
		}

		if damage > 0 {
			res[looter] = e.BonusRewardAmount
		}
	}

	return res
}
