package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

type Enemy struct {
	ID     int
	X      int
	Y      int
	Status bool
}

func NewEnemy(ID, x, y int) Enemy {
	return Enemy{
		ID: ID,
		X:  x,
		Y:  y,
	}
}

type Shot struct {
	TargetID int
	X        int
	Y        int
	Status   bool
}

type TrackStat struct {
	Priority int
	Target   Enemy
}

func NewTrackStat(p int, e Enemy) TrackStat {
	return TrackStat{
		Priority: p,
		Target:   e,
	}
}

func NewShot(tID, x, y int) Shot {
	return Shot{
		TargetID: tID,
		X:        x,
		Y:        y,
	}
}

var TrackerStat = make(map[int]Enemy)

func timestamp() string {
	return time.Now().Format("15:04:05.000")
}

func main() {
	radarChan := make(chan Enemy)
	trackChan := make(chan TrackStat)
	shotChan := make(chan Shot)

	var wg sync.WaitGroup

	fmt.Printf("[%s] [SYSTEM 🔥] Suton is starting...\n", timestamp())

	go radar(radarChan)
	go weapon(trackChan, shotChan)
	go tracker(radarChan, trackChan)
	go stat(shotChan)

	wg.Add(4)
	wg.Wait()

	os.Exit(0)
}

func radar(radarChan chan Enemy) {
	for {
		e := NewEnemy(rand.Intn(100), rand.Intn(50), rand.Intn(80))

		if old, ok := TrackerStat[e.ID]; ok {
			e.X = old.X + rand.Intn(3) - 1
			e.Y = old.Y + rand.Intn(3) - 1
		}

		TrackerStat[e.ID] = e
		fmt.Printf("[%s] [RADAR 📡][INFO] Enemy#%03d located @ (X=%02d, Y=%02d)\n", timestamp(), e.ID, e.X, e.Y)

		radarChan <- e
		time.Sleep(time.Millisecond * 500)
	}
}

func weapon(trackChan chan TrackStat, shotChan chan Shot) {
	for t := range trackChan {
		if t.Priority < 5 {
			fmt.Printf("[%s] [WEAPON 🚀][WARN] Target#%03d has LOW priority (PRI=%d), skipping\n", timestamp(), t.Target.ID, t.Priority)
			continue
		}

		fmt.Printf("[%s] [WEAPON 🚀][FIRE] Engaging Target#%03d (PRI=%d) @ (X=%02d,Y=%02d)\n", timestamp(), t.Target.ID, t.Priority, t.Target.X, t.Target.Y)
		s := NewShot(t.Target.ID, t.Target.X, t.Target.Y)
		shotChan <- s
	}
}

func tracker(radarChan chan Enemy, trackChan chan TrackStat) {
	for e := range radarChan {
		p := rand.Intn(10)
		t := NewTrackStat(p, e)
		TrackerStat[e.ID] = e
		fmt.Printf("[%s] [TRACKER 🎯][INFO] Target#%03d assigned PRI=%d\n", timestamp(), e.ID, p)
		trackChan <- t
	}
}

func stat(shotChan chan Shot) {
	for s := range shotChan {
		e, ok := TrackerStat[s.TargetID]
		if !ok {
			fmt.Printf("[%s] [STAT 📊][ERR] Target#%03d not found in current track records\n", timestamp(), s.TargetID)
			continue
		}
		if (e.X != s.X) || (e.Y != s.Y) {
			fmt.Printf("[%s] [STAT 📊][MISS] Target#%03d evaded! Impact @ (%02d,%02d) ≠ Actual @ (%02d,%02d)\n",
				timestamp(), s.TargetID, s.X, s.Y, e.X, e.Y)
			continue
		}
		fmt.Printf("[%s] [STAT 📊][HIT!!] Target#%03d neutralized @ (X=%02d,Y=%02d)\n", timestamp(), s.TargetID, s.X, s.Y)
		delete(TrackerStat, e.ID)
	}
}
