function test() {
	let r = newKey();
	console.log("newKey", r);
	r = newTxAndSign("2NqXcWAZXfCvkVBZLaFAQ1ksEnF6G4fYRSubmUMckXGG", "HzeXxgjb589tVBs991jAyLUX7wreSZvrWnRxdGQS4co2", "10", "0");
	console.log("newTxAndSign", r);
}
