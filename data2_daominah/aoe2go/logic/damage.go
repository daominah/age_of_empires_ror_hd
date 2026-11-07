package logic

import "math"

// CalculateDamagePerHit is the damage of one normal hit of the attacker on the target, see CalculateShotDamage.
// The attacker may be a building with only its attacks set, e.g. a Castle arrow.
func CalculateDamagePerHit(attacker Unit, target Unit) float64 {
	return CalculateShotDamage(attacker.GetShot(), target)
}

// CalculateShotDamage is the damage of one shot on the target:
// the damage of the main arrow plus the damage of each extra arrow, all assumed to hit.
// Each arrow deals the sum over the armor classes of both of attack minus armor (not below 0), at least 1,
// so against a high armor like a Battering Ram a Maya skirmisher deals 2 instead of 1,
// e.g. a Maya Elite Skirmisher with Hul'che Javelineers throws one javelin of attack 7 and one of attack 1.
func CalculateShotDamage(shot Shot, target Unit) float64 {
	damage := calculateProjectileDamage(shot.Attacks, target)
	extraAttacks := shot.ExtraAttacks
	if extraAttacks == nil {
		extraAttacks = shot.Attacks
	}
	return damage + float64(shot.ExtraProjectiles)*calculateProjectileDamage(extraAttacks, target)
}

// calculateProjectileDamage is the damage of one projectile on the target.
// The target resistance cuts the bonus damage only, e.g. a Halberdier deals 7 + 32 * 0.6 to a Sicilians Light Cavalry.
func calculateProjectileDamage(attacks map[ArmorClass]float64, target Unit) float64 {
	damage := 0.0
	for class, attack := range attacks {
		armor, isFound := target.Armors[class]
		if !isFound {
			continue
		}
		classDamage := math.Max(0, attack-armor)
		if class != ArmorClassMelee && class != ArmorClassPierce {
			classDamage *= 1 - target.BonusDamageResistance
		}
		damage += classDamage
	}
	return math.Max(1, damage)
}

// CalculateHitsToKill is how many hits of the attacker kill the target.
// It uses the HP with decimals, e.g. Mongols Hussar 117.5 HP, shown as 118 in game.
func CalculateHitsToKill(attacker Unit, target Unit) int {
	return int(math.Ceil(cleanHP(target) / CalculateDamagePerHit(attacker, target)))
}

// GetShownHP is the HP shown in game, rounded to nearest (half up),
// e.g. Mongols Hussar 117.5 HP is shown as 118.
func GetShownHP(u Unit) float64 {
	return math.Floor(cleanHP(u) + 0.5)
}

// cleanHP rounds away the tiny error that percent bonuses leave,
// e.g. the Mongols Hussar 117.5 HP comes out of the bonuses as 117.4997.
func cleanHP(u Unit) float64 {
	return math.Round(u.HP*1000) / 1000
}

// CalculateDamagePerTwoSeconds is the average damage of the unit on the target in 2 seconds,
// 2 seconds being the common reload time in the game, so units with different reload times compare fairly.
// It is the damage of CalculateShotDamage, so the armor of the target is included.
// A charged attack replaces a normal shot, assumed to fire as soon as it is charged:
// one charged shot every ChargeTime plus half a reload, the average wait for the next shot,
// e.g. the Elite Bolas Rider throws a strong bolas about every 31 seconds instead of a normal one.
func CalculateDamagePerTwoSeconds(unit Unit, target Unit) float64 {
	if unit.ReloadTime <= 0 {
		return 0
	}
	shotDamage := CalculateShotDamage(unit.GetShot(), target)
	damage := shotDamage * 2 / unit.ReloadTime
	if unit.ChargeAttack != nil {
		cycle := unit.ChargeAttack.ChargeTime + unit.ReloadTime/2
		damage += (CalculateShotDamage(unit.ChargeAttack.Shot, target) - shotDamage) * 2 / cycle
	}
	return damage
}
