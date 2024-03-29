package main

import (
	"encoding/csv"
	"fmt"
	"github.com/fcd/gnfd-payment/monitor"
	"github.com/fcd/gnfd-payment/util"
	"os"
	"sync"
)

func main() {
	gnfdClient := monitor.NewGnfdCompositClients([]string{
		"https://gnfd.qa.bnbchain.world:443",
	}, "greenfield_9000-1741", false)
	gnfdProcessor := monitor.NewGnfdBlockProcessor(gnfdClient)
	start := uint64(3201807)
	end := uint64(3373153)

	each := (end - start) / 10 //171346 blocks in total

	wg := sync.WaitGroup{}
	wg.Add(10)

	for worker := uint64(0); worker < 10; worker++ {
		go func(i uint64) {
			defer wg.Done()
			from := start + i*each
			to := start + (i+1)*each - 1
			if i == 9 {
				to = end
			}
			fmt.Printf("processing blocks from %d to %d, worker %d\n", from, to, i)
			file, err := os.OpenFile(fmt.Sprintf("gnfd_block_processor_%d.csv", i),
				os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				util.Logger.Errorf("fail to open file err: %s", err)
				panic(err)
			}
			defer file.Close()

			writer := csv.NewWriter(file)
			for j := from; j <= to; j++ {
				if j%100 == 0 {
					fmt.Printf("worker %d: processing progress %d, %f \n\n",
						i, j, float64(j-from)/float64(to-from))
				}
				gnfdProcessor.Process(j, writer)
			}
			writer.Write([]string{"done", fmt.Sprintf("%d", i)})
			writer.Flush()

			fmt.Printf("worker %d done\n", i)
		}(worker)
	}

	wg.Wait()
	fmt.Println("done")
}
