package cthulhu

//go:generate mockgen -destination mock/bot.go -package mock -mock_names Service=BotService cthulhu/bot Service
//go:generate mockgen -destination mock/store.go -package mock -mock_names Service=StoreService cthulhu/store Service
