// Package vreme omogoča beleženje temperaturnih meritev vremenskih postaj in izpis njihove statistike.
// Postaje hranimo v slovarju po njihovih oznakah.
// Nov slovar ustvarimo z make(map[string]vreme.Postaja), postaje pa dodamo s funkcijo DodajPostajo.
package vreme

import (
	"fmt"
	"sort"
)

// Postaja predstavlja vremensko postajo z imenom, krajem in meritvami v °C.
type Postaja struct {
	ime     string
	kraj    string
	meritve []float64
}

// DodajPostajo doda novo postajo z dano oznako, imenom in krajem v postaje.
func DodajPostajo(postaje map[string]Postaja, oznaka, ime, kraj string) {
	if _, exists := postaje[oznaka]; exists { // preverimo, ce postaja ze obstaja v slovarju
		fmt.Printf("Postaja %s že obstaja v zbirki.\n", oznaka)
		return
	}
	postaje[oznaka] = Postaja{ime: ime, kraj: kraj} // dodamo novo postajo v slovar
}

// DodajMeritev postaji z dano oznako doda temperaturo temp v °C.
// Temperatura mora biti končno število v dovoljenem razponu [minTemp, maxTemp].
// Če postaje ni ali meritev oziroma razpon ni veljaven, izpiše napako in ne doda meritve.
func DodajMeritev(postaje map[string]Postaja, oznaka string, temp, minTemp, maxTemp float64) {
	postaja, exists := postaje[oznaka]
	if !exists {
		fmt.Printf("Postaje %s ni v zbirki.\n", oznaka)
		return
	}
	if minTemp > maxTemp {
		fmt.Println("Neveljaven temperaturni razpon.")
		return
	}
	if temp < minTemp || temp > maxTemp {
		fmt.Printf("Neveljavna meritev %.1f °C: dovoljen razpon je od %.1f do %.1f °C.\n", temp, minTemp, maxTemp)
		return
	}

	postaja.meritve = append(postaja.meritve, temp)
	postaje[oznaka] = postaja //slovar vsebuje strukture po vrednosti, zato spremenjeno kopijo shranimo nazaj
}

// statistika vrne najmanjšo, največjo in povprečno temperaturo postaje.
// Če postaje ni, stMeritev ni pozitiven ali ni dovolj meritev, vrne ok == false.
func statistika(postaje map[string]Postaja, oznaka string, stMeritev int) (min, max, avg float64, ok bool) {
	postaja, exists := postaje[oznaka]
	if !exists || stMeritev < 1 || len(postaja.meritve) < stMeritev {
		return 0, 0, 0, false
	}

	min = postaja.meritve[0]
	max = postaja.meritve[0]
	vsota := 0.0

	for i := 0; i < len(postaja.meritve); i++ {
		temp := postaja.meritve[i]

		if temp < min {
			min = temp
		}
		if temp > max {
			max = temp
		}
		vsota += temp
	}

	st_m := len(postaja.meritve)
	avg = vsota / float64(st_m)

	return min, max, avg, true
}

// IzpisMeritev izpiše oznake, imena, kraje in meritve vseh postaj.
// Postaje uredi naraščajoče po oznaki, da je vrstni red izpisa vedno enak.
func IzpisMeritev(postaje map[string]Postaja) {
	fmt.Println("MERITVE:")
	for _, oznaka := range urejeneOznake(postaje) {
		postaja := postaje[oznaka]
		fmt.Printf("%s - %s (%s): %v\n", oznaka, postaja.ime, postaja.kraj, postaja.meritve)
	}
}

// IzpisPorocila izpiše ime, najmanjšo, največjo in povprečno temperaturo vsake postaje ter opis Vroče (>= 25 °C), Toplo (>= 15 °C), Hladno (>= 5 °C) ali Mrzlo (< 5 °C).
// Za statistiko zahteva vsaj stMeritev meritev in uporabi skrito funkcijo statistika.
// Če izračun ni mogoč, izpiše »premalo meritev«.
func IzpisPorocila(postaje map[string]Postaja, stMeritev int) {
	fmt.Println("POROČILO:")
	for _, oznaka := range urejeneOznake(postaje) {
		postaja := postaje[oznaka]
		min, max, avg, ok := statistika(postaje, oznaka, stMeritev)
		if !ok {
			fmt.Printf("%s: premalo meritev\n", postaja.ime)
			continue
		}

		var status string

		switch {
		case avg >= 25:
			status = "Vroče"
		case avg >= 15 && avg < 25:
			status = "Toplo"
		case avg >= 5 && avg < 15:
			status = "Hladno"
		default:
			status = "Mrzlo"
		}

		fmt.Printf("%s: min %.1f °C, max %.1f °C, povprečje %.1f °C -> %s\n", postaja.ime, min, max, avg, status)
	}
}

func urejeneOznake(postaje map[string]Postaja) []string {
	oznake := make([]string, 0, len(postaje))
	for oznaka := range postaje {
		oznake = append(oznake, oznaka) //dodamo oznake v seznam, da jih lahko sortiramo
	}
	sort.Strings(oznake) //sortira oznake narascajoce, da je vrstni red izpisa vedno enak
	return oznake
}
