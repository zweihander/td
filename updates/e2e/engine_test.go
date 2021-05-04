package e2e

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
	"github.com/gotd/td/updates"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"golang.org/x/sync/errgroup"
)

func TestEngine(t *testing.T) {
	log := zaptest.NewLogger(t)
	log, _ = zap.NewDevelopment()
	defer log.Sync()

	s := NewSuite()
	s.CommonDialog(
		"pushkin", "Я помню чудное мгновенье:",
		"pushkin", "Передо мной явилась ты,",
		"pushkin", "Как мимолетное виденье,",
		"pushkin", "Как гений чистой красоты.",

		"pushkin", "В томленьях грусти безнадежной,",
		"pushkin", "В тревогах шумной суеты,",
		"pushkin", "Звучал мне долго голос нежный",
		"pushkin", "И снились милые черты.",

		"pushkin", "Шли годы. Бурь порыв мятежный",
		"pushkin", "Рассеял прежние мечты,",
		"pushkin", "И я забыл твой голос нежный,",
		"pushkin", "Твои небесные черты.",

		"pushkin", "В глуши, во мраке заточенья",
		"pushkin", "Тянулись тихо дни мои",
		"pushkin", "Без божества, без вдохновенья,",
		"pushkin", "Без слез, без жизни, без любви.",

		"pushkin", "Душе настало пробужденье:",
		"pushkin", "И вот опять явилась ты,",
		"pushkin", "Как мимолетное виденье,",
		"pushkin", "Как гений чистой красоты.",

		"pushkin", "И сердце бьется в упоенье,",
		"pushkin", "И для него воскресли вновь",
		"pushkin", "И божество, и вдохновенье,",
		"pushkin", "И жизнь, и слезы, и любовь.",
	)

	s.ChannelDialog("gotd_ru",
		"pushkin", "Во глубине сибирских руд",
		"pushkin", "Храните гордое терпенье,",
		"pushkin", "Не пропадет ваш скорбный труд",
		"pushkin", "И дум высокое стремленье.",

		"pushkin", "Несчастью верная сестра,",
		"pushkin", "Надежда в мрачном подземелье",
		"pushkin", "Разбудит бодрость и веселье,",
		"pushkin", "Придет желанная пора:",

		"pushkin", "Любовь и дружество до вас",
		"pushkin", "Дойдут сквозь мрачные затворы,",
		"pushkin", "Как в ваши каторжные норы",
		"pushkin", "Доходит мой свободный глас.",

		"pushkin", "Оковы тяжкие падут,",
		"pushkin", "Темницы рухнут - и свобода",
		"pushkin", "Вас примет радостно у входа,",
		"pushkin", "И братья меч вам отдадут.",
	)

	s.ChannelDialog("gotd_en",
		"pushkin:", "Вихри снежные крутя;",
		"pushkin:", "То, как зверь, она завоет,",
		"pushkin:", "То заплачет, как дитя,",
		"pushkin:", "То по кровле обветшалой",
		"pushkin:", "Вдруг соломой зашумит,",
		"pushkin:", "То, как путник запоздалый,",
		"pushkin:", "К нам в окошко застучит.",

		"pushkin:", "Наша ветхая лачужка",
		"pushkin:", "И печальна и темна.",
		"pushkin:", "Что же ты, моя старушка,",
		"pushkin:", "Приумолкла у окна?",
		"pushkin:", "Или бури завываньем",
		"pushkin:", "Ты, мой друг, утомлена,",
		"pushkin:", "Или дремлешь под жужжаньем",
		"pushkin:", "Своего веретена?",

		"pushkin:", "Выпьем, добрая подружка",
		"pushkin:", "Бедной юности моей,",
		"pushkin:", "Выпьем с горя; где же кружка?",
		"pushkin:", "Сердцу будет веселей.",
		"pushkin:", "Спой мне песню, как синица",
		"pushkin:", "Тихо за морем жила;",
		"pushkin:", "Спой мне песню, как девица",
		"pushkin:", "За водой поутру шла.",

		"pushkin:", "Буря мглою небо кроет,",
		"pushkin:", "Вихри снежные крутя;",
		"pushkin:", "То, как зверь, она завоет,",
		"pushkin:", "То заплачет, как дитя.",
		"pushkin:", "Выпьем, добрая подружка",
		"pushkin:", "Бедной юности моей,",
		"pushkin:", "Выпьем с горя: где же кружка?",
		"pushkin:", "Сердцу будет веселей.",
	)

	s.EncryptedMessage("gotd_bot", "test1")
	s.EncryptedMessage("gotd_bot", "test2")
	s.EncryptedMessage("gotd_bot", "test3")
	s.EncryptedMessage("gotd_bot", "test4")

	server := NewServer(
		s.common,
		s.secret,
		s.channels,
		log.Named("server"),
	)

	state := NewState()
	state.log = log.Named("state")
	engine := updates.New(updates.Config{
		RawClient:    server,
		Handler:      state,
		Channels:     s.initialChannelPts,
		AccessHasher: updates.AccessHasherFunc(func(channelID int) (int64, error) { return 1337, nil }),
		Logger:       log,
	})

	g, ctx := errgroup.WithContext(context.Background())

	uchan := make(chan tg.UpdatesClass, 10)
	g.Go(func() error {
		defer close(uchan)
		return server.Run(ctx, func(u tg.UpdatesClass) error {
			uchan <- u
			return nil
		})
	})

	for i := 0; i < 4; i++ {
		g.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case u, ok := <-uchan:
					if !ok {
						return nil
					}

					if err := engine.HandleUpdates(u); err != nil {
						return err
					}
				}
			}
		})
	}

	require.NoError(t, g.Wait())
	require.Equal(t, s.common, state.common)
	require.Equal(t, s.channels, state.channels)
	require.Equal(t, s.secret, state.secret)
}
