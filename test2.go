			}
			rn = rand.Intn(1000)
			if rn > 995 {
				f.lockSession()
				if len(f.session) > 0 {
					f.session[0] = 48
				}
				f.unlockSession()
				fmt.Println("Messing with sessionn", f.clientID)
			}
		}
		dones = append(dones, doneChan)
	}
	for i := 0; i <= numTrials; i++ {
		r := <-dones[i]
		if r {
			fmt.Println("Status in done is", r)
			done2Chan <- true
			return
		}
	}
	fmt.Println("Finished sending", numTrials)
	done2Chan <- false
}
