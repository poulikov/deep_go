package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

const (
	// Gold: 32 бита, байт 0-3
	goldOffset = 0
	goldBits   = 32

	// Mana: 10 бит, байт 4 (8 бит) + биты 6-7 байта 5 (2 бита)
	manaOffset = 32
	manaBits   = 10

	// Health: 10 бит, биты 0-5 байта 5 (6 бит) + биты 4-7 байта 6 (4 бита)
	healthOffset = 42
	healthBits   = 10

	// Respect: 4 бита, биты 4-7 байта 7
	respectOffset = 56 + 4
	respectBits   = 4

	// Strength: 4 бита, биты 0-3 байта 7
	strengthOffset = 56 + 0
	strengthBits   = 4

	// Experience: 4 бита, биты 4-7 байта 8
	experienceOffset = 64 + 4
	experienceBits   = 4

	// Level: 4 бита, биты 0-3 байта 8
	levelOffset = 64 + 0
	levelBits   = 4

	// House: 1 бит, бит 7 байта 9
	houseOffset = 72 + 7
	houseBits   = 1

	// Gun: 1 бит, бит 6 байта 9
	gunOffset = 72 + 6
	gunBits   = 1

	// Family: 1 бит, бит 5 байта 9
	familyOffset = 72 + 5
	familyBits   = 1

	// Type: 2 бита, биты 3-4 байта 9
	typeOffset = 72 + 3
	typeBits   = 2
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		copy(person.name[:], name[:42])
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.coords = coords{
			x: int32(x),
			y: int32(y),
			z: int32(z),
		}
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person.attrs[:], uint(gold), goldBits, goldOffset)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person.attrs[:], uint(mana), manaBits, manaOffset)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person.attrs[:], uint(health), healthBits, healthOffset)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person.attrs[:], uint(respect), respectBits, respectOffset)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person.attrs[:], uint(strength), strengthBits, strengthOffset)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person.attrs[:], uint(experience), experienceBits, experienceOffset)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person.attrs[:], uint(level), levelBits, levelOffset)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person.attrs[:], 1, houseBits, houseOffset)
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person.attrs[:], 1, gunBits, gunOffset)
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person.attrs[:], 1, familyBits, familyOffset)
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person.attrs[:], uint(personType), typeBits, typeOffset)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type coords struct {
	x int32
	y int32
	z int32
}
type GamePerson struct {
	coords coords
	name   [42]byte
	/*
		golds      32bits,
		mana       10 bits,
		health     10 bits,
		respect    4 bits,
		power      4 bits,
		experience 4 bits,
		lvl        4 bits,
		house      1 bits,
		weapon     1 bits,
		family     1 bits,
		type       2 bits
	*/
	attrs [10]byte
}

func NewGamePerson(options ...Option) GamePerson {
	p := GamePerson{}
	for _, opt := range options {
		opt(&p)
	}
	return p
}

func (p *GamePerson) Name() string {
	return unsafe.String(&p.name[0], 42)
}

func (p *GamePerson) X() int {
	return int(p.coords.x)
}

func (p *GamePerson) Y() int {
	return int(p.coords.y)
}

func (p *GamePerson) Z() int {
	return int(p.coords.z)
}

func (p *GamePerson) Gold() int {
	return int(getBits(p.attrs[:], goldBits, goldOffset))
}

func (p *GamePerson) Mana() int {
	return int(getBits(p.attrs[:], manaBits, manaOffset))
}

func (p *GamePerson) Health() int {
	return int(getBits(p.attrs[:], healthBits, healthOffset))
}

func (p *GamePerson) Respect() int {
	return int(getBits(p.attrs[:], respectBits, respectOffset))
}

func (p *GamePerson) Strength() int {
	return int(getBits(p.attrs[:], strengthBits, strengthOffset))
}

func (p *GamePerson) Experience() int {
	return int(getBits(p.attrs[:], experienceBits, experienceOffset))
}

func (p *GamePerson) Level() int {
	return int(getBits(p.attrs[:], levelBits, levelOffset))
}

func (p *GamePerson) HasHouse() bool {
	return getBits(p.attrs[:], houseBits, houseOffset) == 1
}

func (p *GamePerson) HasGun() bool {
	return getBits(p.attrs[:], gunBits, gunOffset) == 1
}

func (p *GamePerson) HasFamilty() bool {
	return getBits(p.attrs[:], familyBits, familyOffset) == 1
}

func (p *GamePerson) Type() int {
	return int(getBits(p.attrs[:], typeBits, typeOffset))
}

func setBits(data []byte, value uint, bits, offset int) {
	for i := 0; i < bits; i++ {
		byteIndex := (offset + i) / 8
		bitIndex := (offset + i) % 8

		if byteIndex >= len(data) {
			return
		}

		if (value>>i)&1 == 1 {
			data[byteIndex] |= 1 << bitIndex
		} else {
			data[byteIndex] &^= 1 << bitIndex
		}
	}
}

func getBits(data []byte, bits, offset int) uint {
	var result uint
	for i := 0; i < bits; i++ {
		byteIndex := (offset + i) / 8
		bitIndex := (offset + i) % 8

		if byteIndex >= len(data) {
			break
		}

		if (data[byteIndex]>>bitIndex)&1 == 1 {
			result |= 1 << i
		}
	}
	return result
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
