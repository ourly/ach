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
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ourly/ach"
)

func Example_ccdReadDebit() {
	f, err := os.Open(filepath.Join("testdata", "ccd-debit.ach"))
	if err != nil {
		log.Fatal(err)
	}
	r := ach.NewReader(f)
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

	fmt.Printf("Total Amount Debit: %s", strconv.Itoa(achFile.Control.TotalDebitEntryDollarAmountInFile)+"\n")
	fmt.Printf("Total Amount Credit: %s", strconv.Itoa(achFile.Control.TotalCreditEntryDollarAmountInFile)+"\n")
	fmt.Printf("SEC Code: %s", achFile.Batches[0].GetHeader().StandardEntryClassCode+"\n")
	fmt.Printf("CCD Entry Identification Number: %s", achFile.Batches[0].GetEntries()[0].IdentificationNumber+"\n")
	fmt.Printf("CCD Entry Receiving Company: %s", achFile.Batches[0].GetEntries()[0].IndividualName+"\n")
	fmt.Printf("CCD Entry Trace Number: %s", achFile.Batches[0].GetEntries()[0].TraceNumberField()+"\n")
	fmt.Printf("CCD Fee Identification Number: %s", achFile.Batches[0].GetEntries()[1].IdentificationNumber+"\n")
	fmt.Printf("CCD Fee Receiving Company: %s", achFile.Batches[0].GetEntries()[1].IndividualName+"\n")
	fmt.Printf("CCD Fee Trace Number: %s", achFile.Batches[0].GetEntries()[1].TraceNumberField()+"\n")

	// Output:
	// Total Amount Debit: 500125
	// Total Amount Credit: 0
	// SEC Code: CCD
	// CCD Entry Identification Number: location1234567
	// CCD Entry Receiving Company: Best Co. #123456789012
	// CCD Entry Trace Number: 031300010000001
	// CCD Fee Identification Number: Fee123456789012
	// CCD Fee Receiving Company: Best Co. #123456789012
	// CCD Fee Trace Number: 031300010000002
}
