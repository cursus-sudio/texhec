package game

import (
	"domain/blueprints"
	"domain/common/arrays"
	"domain/common/errdesc"
	"domain/common/models"
	domain "domain/models"
	"domain/vo"

	"github.com/ogiusek/ioc"
)

// SERVICE
type DefaultFraction struct {
	Fraction *blueprints.Fraction
}

// SERVICE
type HubErrors struct {
	GameIsFull,
	AlreadyInGame,
	NotInGame,
	NotAHost error
	// "game is already full"
	// "you are already in game"
	// "you are not in game"
	// "only host can do this"

}

type Hub struct {
	models.ModelBase
	GameName     HubName
	GamePassword vo.Hash

	MaxPlayers           uint
	AllowPlayersToChoose bool
	HostId               models.ModelId
	Options              GameOptions
	Users                []*domain.User
	UserFractions        map[models.ModelId]*blueprints.Fraction
}

func NewHub(c ioc.Dic, host *domain.User, gameName HubName, gamePassword vo.Hash) *Hub {
	defaultFraction := ioc.Get[DefaultFraction](c)
	return &Hub{
		ModelBase:    models.NewBase(c),
		GameName:     gameName,
		GamePassword: gamePassword,
		HostId:       host.Id,
		Options:      NewPrepareOptions(c),
		Users: []*domain.User{
			host,
		},
		UserFractions: map[models.ModelId]*blueprints.Fraction{
			host.Id: defaultFraction.Fraction,
		},
	}
}

func (state *Hub) Valid(c ioc.Dic) []error {
	errors := ioc.Get[HubErrors](c)
	var errs []error
	for _, err := range state.GameName.Valid(c) {
		errs = append(errs, errdesc.ErrPath(err).Property("game_name"))
	}
	for _, err := range state.GamePassword.Valid(c) {
		errs = append(errs, errdesc.ErrPath(err).Property("game_password"))
	}
	if int(state.MaxPlayers) == len(state.Users) {
		errs = append(errs, errdesc.ErrPath(errors.GameIsFull).Property("max_players"))
	}
	return errs
}

func (state *Hub) IsFull() bool {
	return int(state.MaxPlayers) == len(state.Users)
}

func (game *Hub) Join(c ioc.Dic, user *domain.User) error {
	errors := ioc.Get[HubErrors](c)
	userIndex := arrays.FirstIndexWhere(game.Users, func(u *domain.User) bool { return u.Id == user.Id })
	if userIndex != -1 {
		return errors.AlreadyInGame
	}

	if game.IsFull() {
		return errors.GameIsFull
	}

	defaultFraction := ioc.Get[DefaultFraction](c)
	game.Users = append(game.Users, user)
	game.UserFractions[user.Id] = defaultFraction.Fraction

	return nil
}

func (state *Hub) Quit(c ioc.Dic, userId models.ModelId) error {
	errors := ioc.Get[HubErrors](c)
	userIndex := arrays.FirstIndexWhere(state.Users, func(u *domain.User) bool { return u.Id == userId })
	if userIndex == -1 {
		return errors.NotInGame
	}
	delete(state.UserFractions, userId)
	state.Users = append(state.Users[:userIndex], state.Users[userIndex+1:]...)
	return nil
}

func (state *Hub) IsHost(userId models.ModelId) bool {
	return state.HostId == userId
}

func (state *Hub) ChangeOptions(c ioc.Dic, changedByUserId models.ModelId, change func(options *GameOptions)) []error {
	errors := ioc.Get[HubErrors](c)
	if !state.IsHost(changedByUserId) {
		return []error{errors.NotAHost}
	}
	change(&state.Options)
	return state.Valid(c)
}

func (state *Hub) ChangeOrder(c ioc.Dic, changedByUserId models.ModelId, changedUserId models.ModelId, newIndex int) error {
	errors := ioc.Get[HubErrors](c)
	changedUserIndex := arrays.FirstIndexWhere(state.Users, func(u *domain.User) bool { return u.Id == changedUserId })
	if changedUserIndex == -1 {
		return errors.NotInGame
	}
	if !state.IsHost(changedByUserId) && !state.AllowPlayersToChoose {
		return errors.NotAHost
	}
	state.Users = arrays.MoveElement(state.Users, changedUserIndex, newIndex)
	return nil
}

func (state *Hub) ChangeFraction(c ioc.Dic, changedByUserId models.ModelId, changedUserId models.ModelId, fraction *blueprints.Fraction) error {
	errors := ioc.Get[HubErrors](c)
	changedUserIndex := arrays.FirstIndexWhere(state.Users, func(u *domain.User) bool { return u.Id == changedUserId })
	if changedUserIndex == -1 {
		return errors.NotInGame
	}
	if !state.IsHost(changedByUserId) && !state.AllowPlayersToChoose {
		return errors.NotAHost
	}
	state.UserFractions[changedUserId] = fraction
	return nil
}
