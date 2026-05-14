package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/mgiks/gopher-social/internal/store"
)

func Seed(store store.Store, db *sql.DB) {
	ctx := context.Background()

	tx, _ := db.BeginTx(ctx, nil)

	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("error creating user:", err)
			return
		}
	}

	tx.Commit()

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("error creating post:", err)
			return
		}
	}

	comments := generateComments(100, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("error creating comment:", err)
			return
		}
	}

	log.Println("seeding complete")
}

var usernames = []string{
	"anna", "lucas", "mike", "zoe", "lily",
	"jack", "olivia", "dan", "emily", "paul",
	"nina", "alex", "kate", "john", "lena",
	"ryan", "jane", "max", "grace", "will",
	"ava", "matt", "ella", "sam", "sophie",
	"noah", "emma", "josh", "ella", "mason",
	"mia", "chris", "lucy", "jake", "ruby",
	"adam", "eva", "tom", "zoe", "hannah",
	"lucas", "lea", "julia", "leo", "luca",
	"liam", "zara", "mila", "ella", "milo",
	"olga",
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := range num {
		username := usernames[i%len(usernames)] + fmt.Sprintf("%d", i)

		user := &store.User{
			Username: username,
			Email:    username + "@example.com",
			Role: store.Role{
				Name: "user",
			},
		}
		user.Password.Set("")
		users[i] = user
	}

	return users
}

var titles = []string{
	"Silent Storm", "Last Breath", "Echoes of Tomorrow", "Midnight Rush", "Fading Light",
	"Broken Wings", "Shadow's Edge", "Crimson Horizon", "Unspoken Truths", "Frostbite",
	"Echoing Steps", "Chasing Shadows", "Silent Fury", "Veil of Mist", "Rising Tide",
	"Twisted Fate", "Out of Time", "Whispers in the Dark", "Cracked Mirror", "Fallen Stars",
}

var contents = []string{
	"Just had the best coffee of my life. Why does a good cup of coffee make everything seem better?",
	"Finally finished a 5K run! It wasn’t easy, but the feeling afterward is totally worth it.",
	"Anyone else obsessed with the new season of that show? I can’t stop binge-watching!",
	"Got some new plants for my apartment. They’re adding such a cozy vibe.",
	"Just got back from a weekend road trip—feeling so refreshed and ready to take on the week!",
	"Trying to eat healthier this week. So far, so good, but the cravings are real!",
	"Spent the afternoon at the park. There’s something so peaceful about being outside on a sunny day.",
	"Can’t believe how fast this year is flying by. Time really does speed up as you get older.",
	"Had a great chat with an old friend today. It's amazing how catching up can make you feel grounded.",
	"Started a new book today. I forgot how much I love getting lost in a good story.",
	"Found a cool local coffee shop today. Can’t believe I’ve been living here for months and never noticed it.",
	"Found my new favorite playlist. It’s perfect for just zoning out and working!",
	"Who else is counting down the days until the weekend? The workweek is dragging for some reason.",
	"Did anyone else try the new restaurant in town? The food was incredible—highly recommend it!",
	"Picked up a new hobby recently—painting. It’s been super relaxing and a fun way to get creative.",
	"Just tried a new workout today, and I feel amazing. Sometimes a good sweat session is all you need.",
	"Trying to get better at journaling. It’s definitely a struggle, but I’m committed to making it a habit.",
	"Finally organized my closet. It feels like a weight has been lifted off my shoulders.",
	"Booked my first solo trip! Nervous but excited to explore a new place by myself.",
	"Feeling super grateful for today. Sometimes it's the small things that make all the difference.",
}

var tags = []string{
	"selfcare", "adventure", "mindfulness", "goodvibes", "travelgoals",
	"weekendvibes", "healthylifestyle", "coffeeaddict", "plantparent", "fitnessjourney",
	"positivity", "exploremore", "booklover", "motivationmonday", "mentalhealthmatters",
	"newbeginnings", "outdooradventures", "creativity", "gratitude", "hustle",
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)

	for i := range num {
		randTags := []string{}
		for range rand.IntN(5-1) + 1 {
			randTags = append(randTags, tags[rand.IntN(len(tags))])
		}

		posts[i] = &store.Post{
			UserID:  users[rand.IntN(len(users))].ID,
			Title:   titles[rand.IntN(len(titles))],
			Content: contents[rand.IntN(len(contents))],
			Tags:    randTags,
		}
	}

	return posts
}

var comments = []string{
	"This looks amazing! I need to try it out ASAP.",
	"Totally agree with you on this. So true!",
	"Wow, I never thought about it that way before.",
	"This is exactly what I needed to hear today. Thanks for sharing!",
	"I’ve been wanting to do this for a while now. Definitely adding it to my list!",
	"Such a great post! Can’t wait to see more from you.",
	"I love this! It’s so inspiring.",
	"Couldn’t have said it better myself. 100% agree!",
	"I have to try this recipe! It looks so delicious.",
	"This post just made my day. Keep up the great work!",
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)

	for i := range num {
		cms[i] = &store.Comment{
			PostID:  posts[rand.IntN(len(posts))].ID,
			UserID:  users[rand.IntN(len(users))].ID,
			Content: comments[rand.IntN(len(comments))],
		}
	}

	return cms
}
