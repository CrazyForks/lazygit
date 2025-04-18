package reflog

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Checkout = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Checkout a reflog commit as a detached head",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.EmptyCommit("two")
		shell.EmptyCommit("three")
		shell.HardReset("HEAD^^")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().ReflogCommits().
			Focus().
			Lines(
				Contains("reset: moving to HEAD^^").IsSelected(),
				Contains("commit: three"),
				Contains("commit: two"),
				Contains("commit (initial): one"),
			).
			SelectNextItem().
			PressPrimaryAction().
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Contains("Checkout branch or commit")).
					Select(MatchesRegexp("Checkout commit [a-f0-9]+ as detached head")).
					Confirm()
			}).
			TopLines(
				Contains("checkout: moving from master to").IsSelected(),
				Contains("reset: moving to HEAD^^"),
			)

		t.Views().Branches().
			Lines(
				Contains("(HEAD detached at").IsSelected(),
				Contains("master"),
			)

		t.Views().Commits().
			Focus().
			Lines(
				Contains("three").IsSelected(),
				Contains("two"),
				Contains("one"),
			)
	},
})
