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
)

// Example_dneWrite writes a DNR file
func Example_dneWrite() {
	fh := mockFileHeader()

	bh := ach.NewBatchHeader()
	bh.ServiceClassCode = ach.CreditsOnly
	bh.CompanyName = "Name on Account"
	bh.CompanyIdentification = fh.ImmediateOrigin
	bh.StandardEntryClassCode = ach.DNE
	bh.CompanyEntryDescription = "Death"
	bh.EffectiveEntryDate = "190816"
	bh.ODFIIdentification = "23138010"
	bh.OriginatorStatusCode = 2

	entry := ach.NewEntryDetail()
	entry.TransactionCode = ach.CheckingReturnNOCCredit
	entry.SetRDFI("031300012")
	entry.DFIAccountNumber = "744-5678-99"
	entry.Amount = 0
	entry.SetOriginalTraceNumber("031300010000001")
	entry.SetReceivingCompany("Best. #1")
	entry.SetTraceNumber(bh.ODFIIdentification, 1)

	addenda := ach.NewAddenda05()
	addenda.PaymentRelatedInformation = `    DATE OF DEATH*010218*CUSTOMERSSN*#########*AMOUNT*$$$$.cc\`
	entry.AddAddenda05(addenda)
	entry.AddendaRecordIndicator = 1

	// build the batch
	batch := ach.NewBatchDNE(bh)
	batch.AddEntry(entry)
	if err := batch.Create(); err != nil {
		log.Fatalf("Unexpected error building batch: %s\n", err)
	}

	// build the file
	file := ach.NewFile()
	file.SetHeader(fh)
	file.AddBatch(batch)
	if err := file.Create(); err != nil {
		log.Fatalf("Unexpected error building file: %s\n", err)
	}

	fmt.Printf("%s", file.Header.String()+"\n")
	fmt.Printf("%s", file.Batches[0].GetHeader().String()+"\n")
	fmt.Printf("%s", file.Batches[0].GetEntries()[0].String()+"\n")
	fmt.Printf("%s", file.Batches[0].GetEntries()[0].Addenda05[0].String()+"\n")
	fmt.Printf("%s", file.Batches[0].GetControl().String()+"\n")
	fmt.Printf("%s", file.Control.String()+"\n")

	// Output:
	// 101 03130001202313801041908161055A094101Federal Reserve Bank   My Bank Name           12345678
	// 5220Name on Account                     231380104 DNEDeath           190816   2231380100000001
	// 621031300012744-5678-99      0000000000031300010000001Best. #1                1231380100000001
	// 705    DATE OF DEATH*010218*CUSTOMERSSN*#########*AMOUNT*$$$$.cc\                  00010000001
	// 82200000020003130001000000000000000000000000231380104                          231380100000001
	// 9000001000001000000020003130001000000000000000000000000
}
