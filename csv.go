package main

import (
    "encoding/csv"
    "fmt"
    "log"
    "os"

    "bitbucket.org/dchapes/ripple/crypto/rkey"
)

func main() {
    var secret string
    // instead of defining data over and over, let's define it once then accumulate things into it.
    data := [][]string{}
    data = append(data, []string{"secret", "address"})
    for i := 0; i < 2000; i++ {
        var s *rkey.FamilySeed
        s, err := rkey.GenerateSeed()
        if err != nil {
            log.Fatal(err)
        }

        pubkey := s.PrivateGenerator.PublicGenerator.Generate(0)
        addr := pubkey.Address()

        if b, err := s.MarshalText(); err != nil {
            log.Fatal(err)
        } else {
            secret = string(b)
        }
        fmt.Printf("Secret:%s\tAddress:%s\n", secret, addr)
        data = append(data, []string{secret, addr})

    }
    // Instead of writing 2000 csv files over each other, let's write one with
    // 2000 lines instead.
    if err := csvExport(data); err != nil {
        log.Fatal(err)
    }
}

// Changed to csvExport, as it doesn't make much sense to export things from
// package main
func csvExport(data [][]string) error {
    file, err := os.Create("result.csv")
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    for _, value := range data {
        if err := writer.Write(value); err != nil {
            return err // let's return errors if necessary, rather than having a one-size-fits-all error handler
        }
    }
    return nil
}
