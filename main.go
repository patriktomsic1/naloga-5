package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/patriktomsic1/naloga-5/vreme"

	"github.com/urfave/cli/v3"
)

func main() {
	if err := cmd().Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func cmd() *cli.Command { // vrne ukaz za belezenje temperaturnih meritev vremenskih postaj
	return &cli.Command{
		Name:  "domaca5",
		Usage: "Beleženje temperaturnih meritev vremenskih postaj",
		Flags: []cli.Flag{
			&cli.FloatFlag{
				Name:  "minTemp",
				Value: -50,
				Usage: "Najnižja dovoljena temperatura meritve v °C",
			},
			&cli.FloatFlag{
				Name:  "maxTemp",
				Value: 50,
				Usage: "Najvišja dovoljena temperatura meritve v °C",
			},
			&cli.IntFlag{
				Name:  "stMeritev",
				Value: 3,
				Usage: "Najmanjše število meritev za izračun statistike",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 0 { // preverimo, da ni dodatnih nepravilnih argumentov
				return fmt.Errorf("nepravilni dodatni argumenti, za pomoč napiši --help")
			}

			minTemp := cmd.Float("minTemp")   // prebere vrednost stikala --minTemp in jo shrani v spremenljivko minTemp, ce stikala ne podamo dobi privzeto vrednost -50
			maxTemp := cmd.Float("maxTemp")   // prebere vrednost stikala --maxTemp in jo shrani v spremenljivko maxTemp, ce stikala ne podamo dobi privzeto vrednost 50
			stMeritev := cmd.Int("stMeritev") // prebere vrednost stikala --stMeritev in jo shrani v spremenljivko stMeritev, ce stikala ne podamo dobi privzeto vrednost 3

			if minTemp > maxTemp {
				return fmt.Errorf("minTemp ne sme biti večji od maxTemp")
			}
			if stMeritev < 1 {
				return fmt.Errorf("stMeritev mora biti vsaj 1")
			}

			prikaziPrimer(minTemp, maxTemp, stMeritev)
			return nil
		},
	}
}

func prikaziPrimer(minTemp, maxTemp float64, stMeritev int) {
	postaje := make(map[string]vreme.Postaja)
	vreme.DodajPostajo(postaje, "KP01", "Koper Luka", "Koper")
	vreme.DodajPostajo(postaje, "KR01", "Kredarica", "Triglav")
	vreme.DodajPostajo(postaje, "LJ01", "Ljubljana Bežigrad", "Ljubljana")
	vreme.DodajPostajo(postaje, "MB01", "Maribor Tabor", "Maribor")

	//Meritve iz DN1 sem dodal skozi javno funkcijo, da tudi zanje veljajo stikala
	zacetneMeritve := []struct {
		oznaka  string
		meritve []float64
	}{
		{"KP01", []float64{24.5, 27, 26.1, 28.4}},
		{"KR01", []float64{2, -1.5, 3.5}},
		{"LJ01", []float64{18.2, 21.5, 19.8}},
		{"MB01", []float64{12.4}},
	}
	for _, podatki := range zacetneMeritve {
		for _, temp := range podatki.meritve {
			vreme.DodajMeritev(postaje, podatki.oznaka, temp, minTemp, maxTemp)
		}
	}

	vreme.DodajMeritev(postaje, "LJ01", 16.5, minTemp, maxTemp) // dodajanje veljavne meritve
	vreme.DodajMeritev(postaje, "LJ01", 65, minTemp, maxTemp)   // dodajanje neveljavne meritve (previsoka temperatura)
	fmt.Println()
	vreme.DodajMeritev(postaje, "PT33", 20, minTemp, maxTemp) // dodajanje meritve za neobstojeco postajo
	fmt.Println()
	vreme.IzpisMeritev(postaje)
	fmt.Println()
	vreme.IzpisPorocila(postaje, stMeritev) //argument stMeritev je dodan ker je v funkciji statistika potrebno vedeti koliko meritev je potrebno za izracun statistike
}
