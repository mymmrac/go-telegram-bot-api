package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

var (
	myID            = tu.ID(331849104)
	groupID         = tu.ID(-1001516926498)
	channelUsername = tu.Username("@mymmrTest")
	groupUsername   = tu.Username("@botesup")
	userUsername    = tu.Username("@mymmrac")
)

const testCase = 35

func main() {
	ctx := context.Background()
	testToken := os.Getenv("TOKEN")

	bot, err := telego.NewBot(testToken,
		telego.WithDefaultDebugLogger(), telego.WithWarnings())
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = bot.GetMe(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch testCase {
	case 1:
		message := &telego.SendMessageParams{
			ChatID: myID,
			Text:   "Test",
			ReplyMarkup: &telego.ReplyKeyboardMarkup{
				Keyboard: [][]telego.KeyboardButton{
					{
						{
							Text: "1",
						},
						{
							Text: "2",
						},
					},
					{
						{
							Text: "3",
						},
					},
				},
				ResizeKeyboard:        true,
				OneTimeKeyboard:       true,
				InputFieldPlaceholder: "Number?",
			},
		}

		msg, err := bot.SendMessage(ctx, message)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(msg)
	case 2:
		updChan, err := bot.UpdatesViaLongPolling(ctx, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		for upd := range updChan {
			fmt.Println(upd)

			if upd.Message != nil {
				_, err := bot.CopyMessage(ctx, &telego.CopyMessageParams{
					ChatID:     telego.ChatID{ID: upd.Message.Chat.ID},
					FromChatID: telego.ChatID{ID: upd.Message.Chat.ID},
					MessageID:  upd.Message.MessageID,
				})
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	case 3:
		p := &telego.ExportChatInviteLinkParams{ChatID: groupID}
		link, err := bot.ExportChatInviteLink(ctx, p)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(*link)
	case 4:
		p := &telego.SendMediaGroupParams{
			ChatID: myID,
			Media: []telego.InputMedia{
				// &telego.InputMediaDocument{
				//	Type:  "document",
				//	Media: telego.InputFile{File: mustOpen("doc.txt")},
				// },
				&telego.InputMediaPhoto{
					Type:  "photo",
					Media: telego.InputFile{URL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRzJZk-efp0id1yxpUHPYwJ1t8vuAwMI_SXfh77dRFWsg1X1ancplws5_DH_WSJ52MHyH8&usqp=CAU"},
				},
				// telego.InputMediaPhoto{
				//	Type:  "photo",
				//	Media: telego.InputFile{URL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTqSw1_1Ar_u3f2lVhYkhz-R0KaaZtDKwx6Y5H1HGceAmx0sqexKzXkSawLG5PRoRKcy6A&usqp=CAU"},
				// },
				&telego.InputMediaPhoto{
					Type:  "photo",
					Media: telego.InputFile{File: mustOpen("img1.jpg")},
				},
				&telego.InputMediaPhoto{
					Type:  "photo",
					Media: telego.InputFile{File: mustOpen("img2.jpg")},
				},
			},
		}
		msgs, err := bot.SendMediaGroup(ctx, p)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, m := range msgs {
			fmt.Println(m)
		}
	case 5:
		err = bot.SetMyCommands(ctx, &telego.SetMyCommandsParams{
			Commands: []telego.BotCommand{
				{
					Command:     "test",
					Description: "Test OK",
				},
			},
			Scope: &telego.BotCommandScopeAllGroupChats{Type: "all_group_chats"},
		})
		if err != nil {
			fmt.Println(err)
			return
		}
	case 6:
		commands, err := bot.GetMyCommands(ctx, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, c := range commands {
			fmt.Println(c.Command, c.Description)
		}
	case 7:
		updParams := &telego.GetUpdatesParams{
			AllowedUpdates: []string{"chat_member"},
		}
		upd, err := bot.GetUpdates(ctx, updParams)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, u := range upd {
			fmt.Println(u.Message.Chat)
		}
	case 8:
		p := &telego.GetChatAdministratorsParams{ChatID: telego.ChatID{ID: -1001516926498}}
		admins, err := bot.GetChatAdministrators(ctx, p)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, u := range admins {
			switch cm := u.(type) {
			case *telego.ChatMemberAdministrator:
				fmt.Println("admin", cm.User)
			case *telego.ChatMemberOwner:
				fmt.Println("owner", cm.User)
			default:
				fmt.Println(cm.MemberStatus())
			}
		}
	case 9:
		dp := &telego.SendDocumentParams{
			ChatID:   myID,
			Document: telego.InputFile{FileID: "BQACAgIAAxkDAAMmYP_FFDZSpqgMsWpK0GCB3hQaI8MAApUPAALeHgABSHe5TRKuQ2NGIAQ"},
			ReplyMarkup: &telego.InlineKeyboardMarkup{InlineKeyboard: [][]telego.InlineKeyboardButton{
				{
					{
						Text:         "Test",
						CallbackData: "1",
					},
				},
			}},
		}
		msg, err := bot.SendDocument(ctx, dp)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(msg.Document)
	case 10:
		dp := &telego.SendDocumentParams{
			ChatID:   myID,
			Document: telego.InputFile{File: mustOpen("doc.txt")},
			Caption:  "Hello world",
		}
		msg, err := bot.SendDocument(ctx, dp)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(msg.Document)
	case 11:
		photo := &telego.SendPhotoParams{
			ChatID:  channelUsername,
			Photo:   telego.InputFile{File: mustOpen("img1.jpg")},
			Caption: "https://test.ua/test_url",
		}

		msg, err := bot.SendPhoto(ctx, photo)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(msg)
	case 12:
		msg := &telego.SendMessageParams{
			ChatID: channelUsername,
			Text:   "Test msg",
		}
		_, err = bot.SendMessage(ctx, msg)
		if err != nil {
			fmt.Println(err)
			return
		}
	case 13:
		msg := &telego.SendMessageParams{
			ChatID: myID,
			Text: `	case 12:
		msg := &telego.SendMessageParams{
			ChatID: channelUsername,
			Text:   "Test msg",
		}
		_, err = bot.SendMessage(msg)
		if err != nil {
			fmt.Println(err)
			return
		}`,
		}

		msg.Entities = []telego.MessageEntity{
			{
				Type:     telego.EntityTypePre,
				Offset:   0,
				Length:   len(msg.Text),
				Language: "go",
			},
		}

		_, err = bot.SendMessage(ctx, msg)
		if err != nil {
			fmt.Println(err)
			return
		}
	case 14:
		_, err := bot.SendMessage(ctx, tu.Message(groupUsername, "Test 1"))
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = bot.SendMessage(ctx, tu.Message(userUsername, "Test 2"))
		if err != nil {
			fmt.Println(err)
			return
		}
	case 15:
		updates, _ := bot.UpdatesViaLongPolling(ctx, nil)

		bh, _ := th.NewBotHandler(bot, updates)
		defer bh.Stop()

		bh.Handle(func(ctx *th.Context, update telego.Update) error {
			fmt.Println(update.Message.Text)
			return nil
		}, func(ctx context.Context, update telego.Update) bool {
			return update.Message != nil
		})

		bh.Handle(func(ctx *th.Context, update telego.Update) error {
			fmt.Println("====")
			fmt.Println(update.Message.Text)
			fmt.Println("====")
			return nil
		}, func(ctx context.Context, update telego.Update) bool {
			return update.Message != nil && update.Message.Text == "OK"
		})

		err = bh.Start()
		assert(err == nil, err)
	case 16:
		updates, _ := bot.UpdatesViaLongPolling(ctx, nil)

		bh, _ := th.NewBotHandler(bot, updates)

		count := 0

		bh.Handle(func(ctx *th.Context, update telego.Update) error {
			fmt.Println("ZERO")
			_, _ = bot.SendMessage(ctx, tu.Message(tu.ID(update.Message.Chat.ID), fmt.Sprintf("Count is zero")))
			count = 1
			return nil
		}, func(ctx context.Context, update telego.Update) bool {
			return update.Message != nil && count == 0
		})

		bh.Handle(func(ctx *th.Context, update telego.Update) error {
			fmt.Println("ONE")
			_, _ = bot.SendMessage(ctx, tu.Message(tu.ID(update.Message.Chat.ID), fmt.Sprintf("Count is one")))
			count = 2
			return nil
		}, func(ctx context.Context, update telego.Update) bool {
			return update.Message != nil && count == 1
		})

		bh.Handle(func(ctx *th.Context, update telego.Update) error {
			fmt.Println("BIG")
			_, _ = bot.SendMessage(ctx, tu.Message(tu.ID(update.Message.Chat.ID), fmt.Sprintf("Count is big: %d", count)))
			count++
			return nil
		}, func(ctx context.Context, update telego.Update) bool {
			return update.Message != nil && count > 1
		})

		defer bh.Stop()
		err = bh.Start()
		assert(err == nil, err)
	case 17:
		updates, _ := bot.UpdatesViaLongPolling(ctx, nil)

		bh, _ := th.NewBotHandler(bot, updates)

		bh.Handle(func(ctx *th.Context, update telego.Update) error {
			msg := update.Message
			matches := th.CommandRegexp.FindStringSubmatch(msg.Text)
			_, _ = bot.SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), fmt.Sprintf("%+v", matches)))
			return nil
		}, th.AnyCommand())

		bh.Handle(func(ctx *th.Context, update telego.Update) error {
			msg := update.Message
			_, _ = bot.SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), fmt.Sprintf("Whaaat? %s", msg.Text)))
			return nil
		}, th.AnyMessage(), th.Not(th.AnyCommand()))

		defer bh.Stop()
		err = bh.Start()
		assert(err == nil, err)
	case 18:
		updates, err := bot.UpdatesViaLongPolling(ctx, nil)
		assert(err == nil, err)

		bh, err := th.NewBotHandler(bot, updates)
		assert(err == nil, err)

		bh.Handle(func(ctx *th.Context, update telego.Update) error {
			msg := update.Message
			_, _ = bot.SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "Running test"))
			return nil
		}, th.CommandEqualArgv("run", "test"))

		bh.Handle(func(ctx *th.Context, update telego.Update) error {
			msg := update.Message
			_, _ = bot.SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "Running update"))
			return nil
		}, th.CommandEqualArgv("run", "update"))

		bh.Handle(func(ctx *th.Context, update telego.Update) error {
			msg := update.Message
			m := tu.Message(tu.ID(msg.Chat.ID), "Run usage:\n```/run test```\n```/run update```")
			m.ParseMode = telego.ModeMarkdownV2
			_, _ = bot.SendMessage(ctx, m)
			return nil
		}, th.Or(
			th.CommandEqualArgc("run", 0),
			th.CommandEqualArgv("help", "run"),
		))

		bh.Handle(func(ctx *th.Context, update telego.Update) error {
			msg := update.Message
			m := tu.Message(tu.ID(msg.Chat.ID), "Unknown subcommand\nRun usage:\n```/run test```\n```/run update```")
			m.ParseMode = telego.ModeMarkdownV2
			_, _ = bot.SendMessage(ctx, m)
			return nil
		}, th.CommandEqual("run"))

		bh.Handle(func(ctx *th.Context, update telego.Update) error {
			msg := update.Message
			_, _ = bot.SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "Help: /run"))
			return nil
		}, th.CommandEqual("help"))

		bh.Handle(func(ctx *th.Context, update telego.Update) error {
			msg := update.Message
			_, _ = bot.SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "Unknown command, use: /run"))
			return nil
		}, th.AnyCommand())

		defer bh.Stop()
		err = bh.Start()
		assert(err == nil, err)
	case 19:
		updates, _ := bot.UpdatesViaLongPolling(ctx, nil)

		bh, _ := th.NewBotHandler(bot, updates)

		bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
			_, _ = bot.SendMessage(ctx, tu.Message(tu.ID(message.Chat.ID), "Hmm?"))
			return nil
		}, th.TextEqual("Hmm"))

		bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
			_, _ = bot.SendMessage(ctx, tu.Message(tu.ID(message.Chat.ID), "Hello"))
			return nil
		})

		defer bh.Stop()
		err = bh.Start()
		assert(err == nil, err)
	case 20:
		img := tu.File(mustOpen("img1.jpg"))
		img2 := tu.File(mustOpen("img2.jpg"))
		audio := tu.File(mustOpen("kitten.mp3"))
		voice := tu.File(mustOpen("kitten.ogg"))
		doc := tu.File(mustOpen("doc.txt"))
		video := tu.File(mustOpen("sample.mp4"))
		note := tu.File(mustOpen("note.mp4"))
		gif := tu.File(mustOpen("cat.mp4"))

		_, err = bot.SendMessage(ctx, tu.Message(myID, "Test"))
		assert(err == nil, err)

		_, err = bot.SendPhoto(ctx, tu.Photo(myID, img))
		assert(err == nil, err)

		_, err = bot.SendAudio(ctx, tu.Audio(myID, audio))
		assert(err == nil, err)

		_, err = bot.SendDocument(ctx, tu.Document(myID, doc))
		assert(err == nil, err)

		time.Sleep(time.Second * 3)

		_, err = bot.SendVideo(ctx, tu.Video(myID, video))
		assert(err == nil, err)

		_, err = bot.SendAnimation(ctx, tu.Animation(myID, gif))
		assert(err == nil, err)

		_, err = bot.SendVoice(ctx, tu.Voice(myID, voice))
		assert(err == nil, err)

		_, err = bot.SendVideoNote(ctx, tu.VideoNote(myID, note))
		assert(err == nil, err)

		time.Sleep(time.Second * 3)

		img = tu.File(mustOpen("img1.jpg"))
		img2 = tu.File(mustOpen("img2.jpg"))

		_, err = bot.SendMediaGroup(ctx, tu.MediaGroup(myID, tu.MediaPhoto(img), tu.MediaPhoto(img2)))
		assert(err == nil, err)

		_, err = bot.SendLocation(ctx, tu.Location(myID, 42, 24))
		assert(err == nil, err)

		_, err = bot.SendVenue(ctx, tu.Venue(myID, 42, 24, "The Thing", "Things str."))
		assert(err == nil, err)

		_, err = bot.SendContact(ctx, tu.Contact(myID, "+424242", "The 42"))
		assert(err == nil, err)

		time.Sleep(time.Second * 3)

		_, err = bot.SendPoll(ctx, tu.Poll(myID, "42?", tu.PollOption("42"), tu.PollOption("24")))
		assert(err == nil, err)

		_, err = bot.SendDice(ctx, tu.Dice(myID, telego.EmojiBasketball))
		assert(err == nil, err)

		err = bot.SendChatAction(ctx, tu.ChatAction(myID, telego.ChatActionTyping))
		assert(err == nil, err)
	case 21:
		updates, _ := bot.UpdatesViaLongPolling(ctx, nil, telego.WithLongPollingUpdateInterval(time.Second))

		bh, _ := th.NewBotHandler(bot, updates)

		bh.HandleInlineQuery(func(ctx *th.Context, query telego.InlineQuery) error {
			err = bot.AnswerInlineQuery(ctx, &telego.AnswerInlineQueryParams{
				InlineQueryID: query.ID,
				Results: []telego.InlineQueryResult{
					&telego.InlineQueryResultArticle{
						Type:                telego.ResultTypeArticle,
						ID:                  "1",
						Title:               "Hmm",
						InputMessageContent: tu.TextMessage("Hmm"),
						ReplyMarkup: tu.InlineKeyboard(tu.InlineKeyboardRow(
							tu.InlineKeyboardButton("GG?").WithCallbackData("ok"),
						)),
					},
				},
			})
			assert(err == nil, err)
			return nil
		})

		bh.HandleCallbackQuery(func(ctx *th.Context, query telego.CallbackQuery) error {
			_, err = bot.EditMessageText(ctx, &telego.EditMessageTextParams{
				Text:            "GG?",
				InlineMessageID: query.InlineMessageID,
			})
			assert(err == nil, err)

			err = bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
				CallbackQueryID: query.ID,
				Text:            "OK",
			})
			assert(err == nil, err)
			return nil
		})

		defer bh.Stop()
		err = bh.Start()
		assert(err == nil, err)
	case 22:
		updates, _ := bot.UpdatesViaLongPolling(ctx, nil)

		bh, _ := th.NewBotHandler(bot, updates)

		auth := func(ctx context.Context, update telego.Update) bool {
			var userID int64

			if update.Message != nil && update.Message.From != nil {
				userID = update.Message.From.ID
			}

			if update.CallbackQuery != nil {
				userID = update.CallbackQuery.From.ID
			}

			if userID == 0 {
				return false
			}

			if userID == 1234 {
				return true
			}

			return false
		}

		bh.Handle(func(ctx *th.Context, update telego.Update) error {
			// DO AUTHORIZED STUFF...
			return nil
		}, auth)

		bh.Handle(func(ctx *th.Context, update telego.Update) error {
			// DO NOT AUTHORIZED STUFF...
			return nil
		}, th.Not(auth))

		defer bh.Stop()
		err = bh.Start()
		assert(err == nil, err)
	case 23:
		updates, _ := bot.UpdatesViaLongPolling(ctx, nil)

		bh, _ := th.NewBotHandler(bot, updates)

		ok := false
		middleware := func(ctx context.Context, update telego.Update) bool {
			return ok
		}

		bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
			ok = true
			fmt.Println("SET OK")
			return nil
		}, th.CommandEqual("ok"))

		bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
			fmt.Println("OK")
			return nil
		}, middleware)

		defer bh.Stop()
		err = bh.Start()
		assert(err == nil, err)
	case 24:
		mux := http.NewServeMux()

		updates, err := bot.UpdatesViaWebhook(ctx, telego.WebhookHTTPServeMux(mux, "POST /"))
		if err != nil {
			panic(err)
		}

		updates, err = bot.UpdatesViaWebhook(ctx, func(handler telego.WebhookHandler) error {
			mux.HandleFunc("POST /", func(writer http.ResponseWriter, request *http.Request) {
				data, err := io.ReadAll(request.Body)
				if err != nil {
					panic(err)
				}

				err = handler(request.Context(), data)
				if err != nil {
					panic(err)
				}

				writer.WriteHeader(http.StatusOK)
			})
			return nil
		})
		if err != nil {
			panic(err)
		}

		err = http.ListenAndServe(":8080", mux)
		if err != nil {
			panic(err)
		}

		<-updates
	case 27:
		note := tu.File(mustOpen("note.mp4"))

		_, err = bot.SendVideoNote(ctx, tu.VideoNote(myID, note))
		assert(err == nil, err)
	case 28:
		err = bot.DeleteWebhook(ctx, nil)
		fmt.Println(err)
	case 29:
		_, err = bot.SendMessage(ctx,
			tu.Message(myID, "Hmm").
				WithReplyMarkup(
					tu.InlineKeyboard(
						tu.InlineKeyboardRow(
							tu.InlineKeyboardButton("OK").
								WithSwitchInlineQueryCurrentChat(""),
						),
					),
				),
		)
		assert(err == nil, err)
	case 30:
		_, err = bot.SendMessage(ctx, tu.Message(myID, "Reply?").
			WithReplyMarkup(tu.ForceReply().WithInputFieldPlaceholder("GG")))
		assert(err == nil, err)
	case 31:
		updates, _ := bot.UpdatesViaLongPolling(ctx, nil)

		bh, _ := th.NewBotHandler(bot, updates)

		bh.Use(th.PanicRecovery())
		bh.Use(
			func(ctx *th.Context, update telego.Update) error {
				fmt.Println("M 2")
				return ctx.Next(update)
			},
			func(ctx *th.Context, update telego.Update) error {
				fmt.Println("M 3")
				return ctx.Next(update)
			},
		)

		bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
			fmt.Println("REGULAR USER")
			return nil
		})

		adminUserMsg := bh.Group(th.AnyMessage(), func(ctx context.Context, update telego.Update) bool {
			ok := update.Message.From.ID == myID.ID
			fmt.Println("OK", ok)
			return ok
		})
		adminUserMsg.HandleMessage(func(ctx *th.Context, message telego.Message) error {
			fmt.Println("ADMIN USER")
			return nil
		})
		adminUserMsg.Use(func(ctx *th.Context, update telego.Update) error {
			fmt.Println("M 4")
			return ctx.Next(update)
		})

		defer bh.Stop()
		err = bh.Start()
		assert(err == nil, err)
	case 32:
		err = bot.CreateNewStickerSet(ctx, &telego.CreateNewStickerSetParams{
			UserID: myID.ID,
			Name:   "the_test_by_ThenWhyBot",
			Title:  "The Test",
			Stickers: []telego.InputSticker{
				{
					Sticker:   telego.InputFile{File: mustOpen("sticker1.png")},
					Format:    telego.StickerFormatStatic,
					EmojiList: []string{"⚡️"},
				},
			},
			StickerType:     telego.StickerTypeRegular,
			NeedsRepainting: false,
		})
		assert(err == nil, err)
	case 33:
		err = bot.AddStickerToSet(ctx, &telego.AddStickerToSetParams{
			UserID: myID.ID,
			Name:   "the_test_by_ThenWhyBot",
			Sticker: telego.InputSticker{
				Sticker: telego.InputFile{
					File: mustOpen("sticker2.png"),
					// URL: "https://upload.wikimedia.org/wikipedia/commons/6/63/Icon_Bird_512x512.png",
				},
				EmojiList: []string{"🐳"},
			},
		})
		assert(err == nil, err)
	case 34:
		err = bot.SetMyDescription(ctx, &telego.SetMyDescriptionParams{
			Description: "",
		})
		assert(err == nil, err)

		err = bot.SetMyShortDescription(ctx, &telego.SetMyShortDescriptionParams{
			ShortDescription: "",
		})
		assert(err == nil, err)
	case 35:
		updates := make(chan telego.Update, 1)

		bh, _ := th.NewBotHandler(nil, updates)

		test := bh.Group()
		test.Use(func(ctx *th.Context, update telego.Update) error {
			fmt.Println("middleware")
			// _ = ctx.Next(update)
			return ctx.Next(update)
		})
		test.Handle(func(ctx *th.Context, update telego.Update) error {
			panic("not here")
		}, th.None())

		test2 := test.Group()
		test2.Use(func(ctx *th.Context, update telego.Update) error {
			fmt.Println("middleware 2")
			return ctx.Next(update)
		})

		test.Handle(func(ctx *th.Context, update telego.Update) error {
			fmt.Println("handler 3")
			return ctx.Next(update)
		})

		bh.Handle(func(ctx *th.Context, update telego.Update) error {
			fmt.Println("handler")
			return nil
		})

		bh.Handle(func(ctx *th.Context, update telego.Update) error {
			fmt.Println("handler 2")
			return nil
		})

		updates <- telego.Update{}
		_ = bh.Start()
	}
}

func assert(ok bool, args ...any) {
	if !ok {
		fmt.Println(args...)
		os.Exit(1)
	}
}

func mustOpen(filename string) *os.File {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return file
}
