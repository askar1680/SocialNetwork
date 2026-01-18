package db

import (
	"AwesomeProject/internal/store"
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var usernames = []string{
	"pixelFalcon",
	"neonTiger",
	"cloudRanger",
	"byteWizard",
	"shadowNova",
	"lunarEcho",
	"ironComet",
	"frostByte",
	"quantumFox",
	"silentDrift",
	"emberKnight",
	"astroLeaf",
	"voidRunner",
	"solarMint",
	"crystalHawk",
	"staticWave",
	"midnightOrb",
	"cobaltCrow",
	"novaSprout",
	"echoBlaze",
	"riftSeeker",
	"plasmaOtter",
	"urbanMyth",
	"cosmoPine",
	"stormGlyph",
	"pixelNomad",
	"velvetSpark",
	"arcaneLoop",
	"neutronIvy",
	"ghostMarble",
	"atomicBreeze",
	"wildSyntax",
	"lunarPanda",
	"cipherBloom",
	"emberGlitch",
	"frozenPulse",
	"skywardByte",
	"obsidianRay",
	"mythicSnow",
	"orbitFable",
	"staticNimbus",
	"solarWisp",
	"crypticAsh",
	"moonlitHex",
	"turboWillow",
	"cosmicDrift",
	"pixelHarbor",
	"auroraKnack",
}

var titles = []string{
	"The Last Signal Before Dawn",
	"Echoes in the Static",
	"A Map Made of Ash and Light",
	"When the Stars Forgot Our Names",
	"The Quiet Between Falling Hours",
	"Notes from a Broken Compass",
	"The City That Dreamed in Color",
	"Shadows Don’t Ask for Permission",
	"Beneath the Neon Weather",
	"The Shape of Tomorrow’s Silence",
	"Messages Left in Orbit",
	"How the Ocean Learned to Listen",
	"Footsteps Along the Vanishing Line",
	"A Theory of Almost Everything",
	"The Night We Borrowed the Sky",
	"Static in the Bloodstream",
	"The Day Time Looked Away",
	"Lanterns for Unfinished Roads",
	"Where the Signal Ends",
	"The Future Is Written in Pencil",
}

var contents = []string{
	"Rain tapped the window while the server logs scrolled endlessly, each line a quiet promise of progress.",
	"A forgotten notebook surfaced in the drawer, filled with half-ideas, arrows, and hope written in pencil.",
	"The train arrived early, carrying strangers who felt oddly familiar in the soft morning light.",
	"Coffee cooled beside the keyboard as thoughts raced faster than fingers could follow.",
	"Somewhere between midnight and sunrise, the bug fixed itself—or maybe patience finally won.",
	"The streetlights flickered like they were thinking, unsure whether the night was truly over.",
	"A single notification changed the mood of the entire afternoon.",
	"The map was outdated, but curiosity insisted on following it anyway.",
	"Music leaked from another room, blending with the hum of machines and distant traffic.",
	"It wasn’t perfect, but it worked—and for now, that was enough.",
}

var tags = []string{
	"technology",
	"creativity",
	"opensource",
	"productivity",
	"storytelling",
	"innovation",
	"minimalism",
	"exploration",
	"automation",
	"design",
	"thinking",
}

var randomComments = []string{
	"Looks solid to me, nice work!",
	"I like the direction this is going.",
	"Small detail, but it makes a big difference.",
	"This feels clean and easy to follow.",
	"Did you consider edge cases here?",
	"Works as expected on my end.",
	"Simple, effective, and readable.",
	"This part could use a short comment.",
	"Nice improvement over the previous version.",
	"Unexpected, but in a good way.",
	"I had the same idea earlier!",
	"Everything flows nicely here.",
	"Good balance between clarity and flexibility.",
	"This solves the problem neatly.",
	"Minor tweak suggestion, otherwise great.",
	"Easy to maintain and extend.",
	"I learned something from this.",
	"Solid choice, no complaints.",
	"This is surprisingly elegant.",
	"Thanks for sharing this!",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()
	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)
	for _, user := range users {
		_ = store.Users.Create(ctx, tx, user)
	}
	//
	posts := generatePosts(200, users)
	for _, post := range posts {
		_ = store.Posts.Create(ctx, post)
	}
	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		_ = store.Comments.CreateComments(ctx, comment)
	}
	fmt.Println("Seeding data is complete")
}

func generateUsers(count int) []*store.User {
	users := make([]*store.User, count)
	for i := 0; i < count; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + strconv.Itoa(i),
			Email:    usernames[i%len(usernames)] + strconv.Itoa(i) + "@gmail.com",
		}
	}
	return users
}

func generatePosts(count int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, count)
	for i := 0; i < count; i++ {
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[i%len(titles)] + strconv.Itoa(i),
			Content: contents[i%len(contents)] + strconv.Itoa(i),
			Tags:    randomTags(),
		}
	}
	return posts
}

func randomTags() []string {
	if len(tags) == 0 {
		return []string{}
	}

	rand.Seed(time.Now().UnixNano())

	// decide how many tags to return (1 to 3)
	n := rand.Intn(3) + 1
	if n > len(tags) {
		n = len(tags)
	}

	// shuffle copy of tags to avoid mutating input
	shuffled := make([]string, len(tags))
	copy(shuffled, tags)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled[:n]
}

func generateComments(count int, users []*store.User, posts []*store.Post) []*store.Comment {
	comments := make([]*store.Comment, count)
	for i := 0; i < count; i++ {
		comments[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: randomComments[rand.Intn(len(randomComments))],
		}
	}
	return comments
}
