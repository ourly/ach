// Licensed to The Moov Authors under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. The Moov Authors licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package examples

import (
	"fmt"
	"github.com/ourly/ach"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func Example_ppdReadSegmentFile() {
	// open a file for reading. Any io.Reader Can be used
	fCredit, err := os.Open(filepath.Join("testdata", "segmentFile-ppd-credit.ach"))
	if err != nil {
		log.Fatal(err)
	}
	r := ach.NewReader(fCredit)
	achFile, err := r.Read()
	if err != nil {
		fmt.Printf("Issue reading file: %+v \n", err)
	}
	// ensure we have a validated file structure
	if achFile.Validate(); err != nil {
		fmt.Printf("Could not validate entire read file: %v", err)
	}
	// If you trust the file but it's formatting is off building will probably resolve the malformed file.
	if err := achFile.Create(); err != nil {
		fmt.Printf("Could not create file with read properties: %v", err)
	}

	// open a file for reading. Any io.Reader Can be used
	fDebit, err := os.Open(filepath.Join("testdata", "segmentFile-ppd-debit.ach"))
	if err != nil {
		log.Fatal(err)
	}
	rDebit := ach.NewReader(fDebit)
	achFileDebit, err := rDebit.Read()
	if err != nil {
		fmt.Printf("Issue reading file: %+v \n", err)
	}
	// ensure we have a validated file structure
	if achFileDebit.Validate(); err != nil {
		fmt.Printf("Could not validate entire read file: %v", err)
	}
	// If you trust the file but it's formatting is off building will probably resolve the malformed file.
	if achFileDebit.Create(); err != nil {
		fmt.Printf("Could not create file with read properties: %v", err)
	}

	fmt.Printf("Total Credit Amount: %s", strconv.Itoa(achFile.Control.TotalCreditEntryDollarAmountInFile)+"\n")
	fmt.Printf("SEC Code: %s", achFile.Batches[0].GetHeader().StandardEntryClassCode+"\n")
	fmt.Printf("Total Debit Amount: %s", strconv.Itoa(achFileDebit.Control.TotalDebitEntryDollarAmountInFile)+"\n")
	fmt.Printf("SEC Code: %s", achFileDebit.Batches[0].GetHeader().StandardEntryClassCode+"\n")

	// Output:
	// Total Credit Amount: 200000000
	// SEC Code: PPD
	// Total Debit Amount: 200000000
	// SEC Code: PPD
}
