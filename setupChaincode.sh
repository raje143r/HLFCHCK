
echo "Installing tfbc chaincode to peer0.buyerbank.tfbc.com..."

# install chaincode
# Install code on bank peer
docker exec -e "CORE_PEER_LOCALMSPID=BuyerBankMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/buyerbank.tfbc.com/users/Admin@buyerbank.tfbc.com/msp" -e "CORE_PEER_ADDRESS=peer0.buyerbank.tfbc.com:7051" cli peer chaincode install -n tfbccc -v 1.0 -p github.com/tfbc/go -l golang

echo "Installed tfbc chaincode to peer0.buyerbank.tfbc.com"

echo "Installing tfbc chaincode to peer0.sellerbank.tfbc.com..."
docker exec -e "CORE_PEER_LOCALMSPID=SellerBankMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/sellerbank.tfbc.com/users/Admin@sellerbank.tfbc.com/msp" -e "CORE_PEER_ADDRESS=peer0.sellerbank.tfbc.com:7051" cli peer chaincode install -n tfbccc -v 1.0 -p github.com/tfbc/go -l golang
echo "Installed tfbc chaincode to peer0.sellerbank.tfbc.com"


echo "Installing tfbc chaincode to peer0.carrier.tfbc.com"
docker exec -e "CORE_PEER_LOCALMSPID=CarrierMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/carrier.tfbc.com/users/Admin@carrier.tfbc.com/msp" -e "CORE_PEER_ADDRESS=peer0.carrier.tfbc.com:7051" cli peer chaincode install -n tfbccc -v 1.0 -p github.com/tfbc/go -l golang
echo "Installed tfbc chaincode to peer0.carrier.tfbc.com"



echo "Installing tfbc chaincode to peer0.buyer.tfbc.com...."

# Install code on buyer peer
docker exec -e "CORE_PEER_LOCALMSPID=BuyerMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/buyer.tfbc.com/users/Admin@buyer.tfbc.com/msp" -e "CORE_PEER_ADDRESS=peer0.buyer.tfbc.com:7051" cli peer chaincode install -n tfbccc -v 1.0 -p github.com/tfbc/go -l golang

echo "Installed tfbc chaincode to peer0.buyer.tfbc.com"

echo "Installing tfbc chaincode to peer0.seller.tfbc.com..."
# Install code on seller peer
docker exec -e "CORE_PEER_LOCALMSPID=SellerMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/seller.tfbc.com/users/Admin@seller.tfbc.com/msp" -e "CORE_PEER_ADDRESS=peer0.seller.tfbc.com:7051" cli peer chaincode install -n tfbccc -v 1.0 -p github.com/tfbc/go -l golang

sleep 5

echo "Installed tfbc chaincode to peer0.seller.tfbc.com"

echo "Instantiating tfbc chaincode.."

docker exec -e "CORE_PEER_LOCALMSPID=BuyerBankMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/buyerbank.tfbc.com/users/Admin@buyerbank.tfbc.com/msp" -e "CORE_PEER_ADDRESS=peer0.buyerbank.tfbc.com:7051" cli peer chaincode instantiate -o orderer.tfbc.com:7050 -C tfbcchannel -n tfbccc -l golang -v 1.0 -c '{"Args":[""]}' -P "OR ('BuyerBankMSP.member','BuyerMSP.member','SellerMSP.member','SellerBankMSP.member','CarrierMSP.member')"

echo "Instantiated tfbc chaincode."

echo "Following is the docker network....."

docker ps
