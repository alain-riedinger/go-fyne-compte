// [Fyne toolkit documentation for developers | Develop using Fyne](https://developer.fyne.io/index.html)

package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Create the application and window
	myApp := app.New()
	myWindow := myApp.NewWindow("Le Compte est Bon")

	// Global string for textual solution
	var solToDisplay *Solution
	const solEmpty = "\n\n\n\n\n\n\n\n\n\n\n\n"

	// Create the items of the window
	const nbPlaques = 6
	var btnPlaques [nbPlaques](*widget.Button)
	for i := 0; i < len(btnPlaques); i++ {
		btnPlaques[i] = widget.NewButton("___", tapped)
	}
	// Items are displayed horizontally in a grid
	gridPlaques := container.New(layout.NewGridLayout(6),
		btnPlaques[0], btnPlaques[1], btnPlaques[2], btnPlaques[3], btnPlaques[4], btnPlaques[5])

	btnTirage := widget.NewButton("___", tapped)
	progress := widget.NewProgressBar()
	// Button is on the left, with its default size
	// Progress bar is on the left, stretched to use all the remaining space
	gridTirage := container.New(layout.NewBorderLayout(nil, nil, btnTirage, nil), btnTirage, progress)

	// Separator makes the layout nice
	separator1 := widget.NewSeparator()

	// Trigger for solution display
	// Text is forced on several lines
	txtSolution := widget.NewTextGridFromString(solEmpty)
	btnSolution := widget.NewButton("Solution?", func() {
		// Output final result
		state := "Exact"
		if solToDisplay.Best.Value != solToDisplay.Tirage {
			state = "Approched"
		}
		nbReturns := strings.Count(solToDisplay.Best.Text, "\n")
		txtSolution.SetText(fmt.Sprintf("Solution [%s]\n\n%s%s",
			state,
			solToDisplay.Best.Text,
			strings.Repeat("\n", 10-nbReturns)))
	})

	// Separator makes the layout nice
	separator2 := widget.NewSeparator()

	// Channel for stopping the time of the progress bar
	// stop := make(chan bool)

	// Buttons with actions
	newGame := widget.NewButton("Play!", func() {
		cpt := NewCompte()

		txtSolution.SetText(solEmpty)

		plaques := cpt.GetPlaques()
		for i := 0; i < len(btnPlaques); i++ {
			btnPlaques[i].SetText(strconv.Itoa(plaques[i]))
		}

		tirage := cpt.GetTirage()
		btnTirage.SetText(strconv.Itoa(tirage))

		// go timer(stop, progress)
		countup(progress)

		// Solution is searched during chrono time (it's shorter so no cheating)
		sol := make(chan *Solution)
		go findSolution(cpt, tirage, plaques, sol)
		solToDisplay = <-sol
		close(sol)
	})
	quit := widget.NewButton("Quit", func() {
		myApp.Quit()
	})
	gridActions := container.New(layout.NewGridLayout(2), newGame, quit)

	// Compose the window with the items
	// Items are horizontally stacked, first parameter is on top, and so on downwards
	myWindow.SetContent(container.New(layout.NewVBoxLayout(),
		gridPlaques,
		gridTirage,
		separator1,
		btnSolution,
		txtSolution,
		separator2,
		gridActions))

	// Trigger the progress bar
	progress.Min = 0
	progress.Max = 40
	progress.TextFormatter = func() string {
		// No percent displayed to avoid distraction
		return ""
	}

	// Master loop that runs the widow
	myWindow.ShowAndRun()
}

func timer(stop chan bool, pb *widget.ProgressBar) {
	for t := 0.0; t <= pb.Max; t++ {
		select {
		case <-stop:
			return
		default:
			time.Sleep(1 * time.Second)
			pb.SetValue(t)
		}
	}
}

// tapped is a dummy function, for the plaques and tirage to be buttons
func tapped() {
	// Nothing to do
}

func findSolution(cpt *Compte, tirage int, plaques []int, sol chan *Solution) {
	// Initialize the recursive calculation root structure
	solution := NewSolution()
	solution.Tirage = tirage
	solution.Depth = len(plaques)
	res := NewResult()
	res.Steps = plaques
	res.Text = ""
	res.Value = 0
	sort.Ints(res.Steps)
	solution.Current = append(solution.Current, res)

	// Initialize the best approaching structure
	solution.Best = NewResult()
	solution.Best.Steps = res.Steps
	solution.Best.Value = solution.Best.Steps[len(solution.Best.Steps)-1]
	solution.Best.Text = fmt.Sprintf("%d", solution.Best.Value)

	// Start the recursive resolution
	solution = cpt.SolveTirage(*solution)

	// Send solution to the channel, while execution in parallel
	sol <- solution
}

func countup(pb *widget.ProgressBar) {
	up := 0.0
	timer := time.Tick(1 * time.Second)
	for up < pb.Max {
		<-timer
		pb.SetValue(up)
		up += 1.0
	}
}
