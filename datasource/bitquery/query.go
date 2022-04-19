package bitquery

const queryTokenInfo = `query (
  $network:  EthereumNetwork!,
   $address: String!, 
  $from: ISO8601DateTime, 
  $till: ISO8601DateTime
){
  ethereum(network: $network) {
   smartContractCalls(date: {since: $from, till: $till},
    smartContractAddress: {is: $address})
  	{
      created: minimum(of: block, get: time)
      crated_block: minimum(of: block)
      created_tx: minimum(of: block, get: tx_hash)
      created_by: minimum(of: block, get: caller)
    }
  }
}`
