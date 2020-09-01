# A Wallet #
A simple wallet for doing (very) basic operations on the Stellar network:.

## Basic Account Management ##
Create accounts:

    ./wallet  
    [*] Generating address ...
      public key: GAF7V7IVFVOEAUCHHLV5OS4QAGPCJB7KMLL5FAF63VJW4BZCX5XDBYFJ
      private key: SBBO6KETVVHI2AY3CDVM4NSI7XHXJKLPJTKMVSGWE2NSN6JVIY5WQOHV
    [*] Saving to ./accounts/account.json ...
    [*] Funding GAF7V7IVFVOEAUCHHLV5OS4QAGPCJB7KMLL5FAF63VJW4BZCX5XDBYFJ ...
      Hash: 3c62ce4abfd1f330aebdd72e3db0afb46fe6c1c09ca34bb81d2484081faa73fb
    [*] Retrieving balances for GAF7V7IVFVOEAUCHHLV5OS4QAGPCJB7KMLL5FAF63VJW4BZCX5XDBYFJ ...
      10000.0000000 (XLM)

Load / check accounts:

     ./wallet -load accounts/account.json 
    [*] Loading account from accounts/account.json ...
    [*] Retrieving balances for GAF7V7IVFVOEAUCHHLV5OS4QAGPCJB7KMLL5FAF63VJW4BZCX5XDBYFJ ...
      10000.0000000 (XLM)

## Making Payments ##
Transfer lumens between accounts:

    ./wallet -load accounts/account.json -dest accounts/account-1.json -amount 1234
    [*] Loading account from accounts/account.json ...
    [*] Retrieving balances for GAF7V7IVFVOEAUCHHLV5OS4QAGPCJB7KMLL5FAF63VJW4BZCX5XDBYFJ ...
      10000.0000000 (XLM)
    [*] Loading account from accounts/account-1.json ...
    [*] Sending 1234 XLM to GAA76GBBRN2VZEG5VHGDMIYM3QHXJD3HTCMDTEKZL4NHOS3XFXEGGINZ ...
      Ledger: 469874
      Hash: 524fad15cd02362f9800f9c9b4f081912659b79ce4ebc90ae1dd6829a711d831

Issue tokens (create trustlines):

    [*] Loading account from accounts/account.json ...
    [*] Retrieving balances for GAF7V7IVFVOEAUCHHLV5OS4QAGPCJB7KMLL5FAF63VJW4BZCX5XDBYFJ ...
      8765.9999900 (XLM)
    [*] Loading account from accounts/account-1.json ...
    [*] Opening trustline to GAA76GBBRN2VZEG5VHGDMIYM3QHXJD3HTCMDTEKZL4NHOS3XFXEGGINZ for DOGE ...
      Ledger: 469880
      Hash: a95b81985e99c3b8cfab287c00969bfcd808325073e2b7ab46a1c337092bbcfd
    [*] Sending 1234 DOGE to GAA76GBBRN2VZEG5VHGDMIYM3QHXJD3HTCMDTEKZL4NHOS3XFXEGGINZ ...
      Ledger: 469881
      Hash: 70c7c26a5ea26971a7eb10eb45501e92ec896e4db5463ca091be02875a88281d

Transferring the issued token won't recreate the trustline:

    ./wallet -load accounts/account.json -dest accounts/account-1.json -amount 42 -asset DOGE  
    [*] Loading account from accounts/account.json ...
    [*] Retrieving balances for GAF7V7IVFVOEAUCHHLV5OS4QAGPCJB7KMLL5FAF63VJW4BZCX5XDBYFJ ...
      8765.9999700 (XLM)
    [*] Loading account from accounts/account-1.json ...
    [*] Opening trustline to GAA76GBBRN2VZEG5VHGDMIYM3QHXJD3HTCMDTEKZL4NHOS3XFXEGGINZ for DOGE ...
      Skipping: already exists.
    [*] Sending 42 DOGE to GAA76GBBRN2VZEG5VHGDMIYM3QHXJD3HTCMDTEKZL4NHOS3XFXEGGINZ ...
      Ledger: 469886
      Hash: b1534f67fee47296318c47e66ace4de78be9fd3d9be9bf498bc29f8b48dae524

And, obviously, balances update accordingly:

    ./wallet -load accounts/account-1.json
    [*] Loading account from accounts/account-1.json ...
    [*] Retrieving balances for GAA76GBBRN2VZEG5VHGDMIYM3QHXJD3HTCMDTEKZL4NHOS3XFXEGGINZ ...
      1276.0000000 (DOGE)
      14533.9999400 (XLM)