package app

import (
	"fmt"
	"time"

	"github.com/oseducation/knowledge-graph/model"
	"github.com/pkg/errors"
)

const BotID = "aiTutorBotID01234567890123"

func (a *App) AskQuestion(userPost *model.Post) (string, error) {
	posts, err := a.Store.Post().GetPosts(&model.PostGetOptions{
		UserID: userPost.UserID,
	})
	if err != nil {
		return "", errors.Wrapf(err, "userPost = %v", userPost)
	}
	if len(posts) == 0 {
		return "Welcome to the Knowledge Graph!", nil
	}
	return "You have a new post!", nil
}

func (a *App) CreateFirstPost(oldPosts []*model.Post, userID string) (*model.Post, error) {
	user, err := a.Store.User().Get(userID)
	if err != nil {
		return nil, errors.Wrapf(err, "can't get user = %v", userID)
	}
	user.Sanitize()

	message := ""
	var options []model.Option
	if len(oldPosts) == 0 {
		message = fmt.Sprintf(`👋 Hello %s, and welcome to Vitsi AI! I'm your AI tutor, here to help you learn programming and navigate the world of code. 🚀
You might know nothing about coding, but it's OK. You'll start coding in minutes. In fact, our first goal will be to write our first program.
Feel free to ask any questions, and let's make your coding journey exciting and fruitful! 🎉`, user.FirstName)

		options = []model.Option{{
			ID:                "1",
			TextOnButton:      "Let's start!",
			MessageAfterClick: "Great, Let's start with the next topic!",
			Action:            model.PostActionTypeNextTopic,
		}}
	} else {
		lastPostDate := oldPosts[len(oldPosts)-1].UpdatedAt
		timePassed := model.GetMillis() - lastPostDate
		if timePassed > 7*24*time.Hour.Milliseconds() { // more than 7 days
			message = fmt.Sprintf(`🎉 Welcome back %s! It’s great to see you again. I hope everything is going well! 🌟
If you’re ready to dive back into learning, I’m here to help you pick up right where you left off or explore new topics. Just let me know what you’d like to focus on today. And remember, no question is too small — I’m here to support you on your coding journey! 💻✨
If there’s anything specific you’d like to review or learn next, please tell me, or we can continue with the next topic in your learning path. Let’s get coding! 🚀`, user.FirstName)
			options = []model.Option{{
				ID:                "1",
				TextOnButton:      "Let's continue!",
				MessageAfterClick: "Let's continue with the next topic!",
				Action:            model.PostActionTypeNextTopic,
			}, {
				ID:                "2",
				TextOnButton:      "Set A Different Goal!",
				MessageAfterClick: "Let's set a different goal!",
				Action:            model.PostActionTypeLink,
				Link:              "/dashboard/goals",
			}}
		} else if timePassed > 24*time.Hour.Milliseconds() { // more than a day
			message = fmt.Sprintf(`👋 Welcome back %s! It’s fantastic to see you continuing your coding journey today. Consistency is key to mastering programming, and you’re doing great! 🎉
Let’s keep the momentum going! If you’re ready to dive back in, I can help you pick up where we left off, or if there’s something specific you’d like to focus on, feel free to let me know. Your progress is impressive, and I’m here to support you every step of the way. 💪
What would you like to work on today? Let’s make the most out of this session! 🚀`, user.FirstName)
			options = []model.Option{{
				ID:                "1",
				TextOnButton:      "Let's continue!",
				MessageAfterClick: "Let's continue with the topic in progress!",
				Action:            model.PostActionTypeNextTopic,
			}}
		} else if timePassed > time.Hour.Milliseconds() { // more than 1 hours
			message = `🚀 Welcome back! It's great to see you so dedicated to learning. Your focus and determination are truly commendable! 👏
You're on the right track, and I'm here to help you keep the momentum going. Whether you'd like to continue with the previous topic, tackle a new challenge, or need clarification on anything, just let me know. Together, we can make every coding session count! 💻✨`
			options = []model.Option{{
				ID:                "1",
				TextOnButton:      "Let's continue!",
				MessageAfterClick: "Let's continue with the topic in progress!",
				Action:            model.PostActionTypeNextTopic,
			}}
		} else if timePassed > 5*time.Minute.Milliseconds() { // more than 5 minutes
			message = `Let's continue with the topics in progress!`
			options = []model.Option{{
				ID:                "1",
				TextOnButton:      "OK!",
				MessageAfterClick: "OK!",
				Action:            model.PostActionTypeNextTopic,
			}}
		} else {
			return nil, nil
		}
	}

	post, err := a.CreatePost(&model.Post{
		LocationID: fmt.Sprintf("%s_%s", userID, BotID),
		UserID:     BotID,
		Message:    message,
		PostType:   model.PostTypeWithActions,
		Props: map[string]interface{}{
			"options": options,
		},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "can't create post")
	}

	return post, nil
}