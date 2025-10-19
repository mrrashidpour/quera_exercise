package main

func StartDecipher(senderChan chan string, decipherer func(encrypted string) string) chan string {
	resultChan := make(chan string, 5)

	go func() {
		for msg := range senderChan {
			deciphered := decipherer(msg)
			resultChan <- deciphered
		}
		close(resultChan)
	}()

	return resultChan
}
