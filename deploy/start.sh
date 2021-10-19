#!/bin/bash

set -e
mkdir -p channel-artifacts

ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/gdzce.cn/orderers/order.gdzce.cn/msp/tlscacerts/tlsca.gdzce.cn-cert.pem
NODE1_ORG1_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization1.gdzce.cn/peers/node1.organization1.gdzce.cn/tls/ca.crt
NODE1_ORG2_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization2.gdzce.cn/peers/node1.organization2.gdzce.cn/tls/ca.crt
NODE1_ORG3_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization3.gdzce.cn/peers/node1.organization3.gdzce.cn/tls/ca.crt
CHAINCODE_NAME=edu-mgmt
CHAINCODE_VERSION=1.0.0

echo "生成证书和起始区块信息"
cryptogen generate --config=./crypto-config.yaml
configtxgen -profile SampleMultiNodeEtcdRaft -channelID syschannel -outputBlock ./channel-artifacts/genesis.block

echo "生成通道的TX文件"
configtxgen -profile ThreeOrgsChannel -outputCreateChannelTx ./channel-artifacts/channelname.tx -channelID channelname

echo "生成锚节点配置更新文件"
configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Organization1MSPanchors.tx -channelID channelname -asOrg Organization1MSP
configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Organization2MSPanchors.tx -channelID channelname -asOrg Organization2MSP
configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Organization3MSPanchors.tx -channelID channelname -asOrg Organization3MSP

echo "区块链网络 ： 启动"
docker-compose up -d
echo "正在等待节点的启动完成，等待15秒"
sleep 15

# 五、在区块链上按照刚刚生成的TX文件去创建通道
# 该操作和上面操作不一样的是，这个操作会写入区块链
echo "五、在区块链上按照刚刚生成的TX文件去创建通道"
docker exec cli peer channel create -o order.gdzce.cn:7050 -c channelname -f ./channel-artifacts/channelname.tx --tls true --cafile $ORDERER_CA


# 六、让节点去加入到通道
# 所有节点都要加入通道中
echo "六、让节点去加入到通道"
docker exec cli peer channel join -b channelname.block
# 修改环境变量链接到其他节点
#docker exec -e "CORE_PEER_TLS_ROOTCERT_FILE="$node1_ORG3_CA -e "CORE_PEER_LOCALMSPID=Organization3MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization3.gdzce.cn/users/Admin@organization3.gdzce.cn/msp" -e "CORE_PEER_ADDRESS=node1.organization3.gdzce.cn:11051" cli peer channel join -b channelname.block
docker exec -e CORE_PEER_TLS_ROOTCERT_FILE=$NODE1_ORG3_CA -e "CORE_PEER_LOCALMSPID=Organization3MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization3.gdzce.cn/users/Admin@organization3.gdzce.cn/msp" -e "CORE_PEER_ADDRESS=node1.organization3.gdzce.cn:11051" cli peer channel join -b channelname.block
docker exec -e CORE_PEER_TLS_ROOTCERT_FILE=$NODE1_ORG3_CA -e "CORE_PEER_LOCALMSPID=Organization3MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization3.gdzce.cn/users/Admin@organization3.gdzce.cn/msp" -e "CORE_PEER_ADDRESS=node2.organization3.gdzce.cn:12051" cli peer channel join -b channelname.block
docker exec -e CORE_PEER_TLS_ROOTCERT_FILE=$NODE1_ORG1_CA -e "CORE_PEER_LOCALMSPID=Organization1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization1.gdzce.cn/users/Admin@organization1.gdzce.cn/msp" -e "CORE_PEER_ADDRESS=node2.organization1.gdzce.cn:8051" cli peer channel join -b channelname.block
docker exec -e CORE_PEER_TLS_ROOTCERT_FILE=$NODE1_ORG2_CA -e "CORE_PEER_LOCALMSPID=Organization2MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization2.gdzce.cn/users/Admin@organization2.gdzce.cn/msp" -e "CORE_PEER_ADDRESS=node1.organization2.gdzce.cn:9051" cli peer channel join -b channelname.block
docker exec -e CORE_PEER_TLS_ROOTCERT_FILE=$NODE1_ORG2_CA -e "CORE_PEER_LOCALMSPID=Organization2MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization2.gdzce.cn/users/Admin@organization2.gdzce.cn/msp" -e "CORE_PEER_ADDRESS=node2.organization2.gdzce.cn:10051" cli peer channel join -b channelname.block

# 六.一 更新锚节点通道
echo "七、更新锚节点到通道"
#docker exec -e "CORE_PEER_TLS_ROOTCERT_FILE="$node1_ORG3_CA -e "CORE_PEER_LOCALMSPID=Organization3MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization3.gdzce.cn/users/Admin@organization3.gdzce.cn/msp" -e "CORE_PEER_ADDRESS=node1.organization3.gdzce.cn:11051" cli peer channel update -o order.gdzce.cn:7050 --tls true --cafile $ORDERER_CA -c channelname -f ./channel-artifacts/Organization3MSPanchors.tx
docker exec -e "CORE_PEER_LOCALMSPID=Organization3MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization3.gdzce.cn/users/Admin@organization3.gdzce.cn/msp" -e "CORE_PEER_ADDRESS=node1.organization3.gdzce.cn:11051" cli peer channel update -o order.gdzce.cn:7050 -c channelname -f ./channel-artifacts/Organization3MSPanchors.tx --tls true --cafile $ORDERER_CA
docker exec -e "CORE_PEER_LOCALMSPID=Organization1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization1.gdzce.cn/users/Admin@organization1.gdzce.cn/msp" -e "CORE_PEER_ADDRESS=node1.organization1.gdzce.cn:7051" cli peer channel update -o order.gdzce.cn:7050 -c channelname -f ./channel-artifacts/Organization1MSPanchors.tx --tls true --cafile $ORDERER_CA
docker exec -e "CORE_PEER_LOCALMSPID=Organization2MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization2.gdzce.cn/users/Admin@organization2.gdzce.cn/msp" -e "CORE_PEER_ADDRESS=node1.organization2.gdzce.cn:9051" cli peer channel update -o order.gdzce.cn:7050 -c channelname -f ./channel-artifacts/Organization2MSPanchors.tx --tls true --cafile $ORDERER_CA

echo -e "====== 安装链码 ${CHAINCODE_NAME} ${CHAINCODE_VERSION} ======"
echo "node1.organization1 安装链码"
docker exec -e CORE_PEER_TLS_ROOTCERT_FILE=$NODE1_ORG1_CA -e CORE_PEER_ADDRESS=node1.organization1.gdzce.cn:7051 -e CORE_PEER_LOCALMSPID=Organization1MSP -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization1.gdzce.cn/users/Admin@organization1.gdzce.cn/msp cli peer chaincode install -n ${CHAINCODE_NAME} -v ${CHAINCODE_VERSION} -l golang -p github.com/chaincode/${CHAINCODE_NAME}

echo "node1.organization3 安装链码"
docker exec -e CORE_PEER_TLS_ROOTCERT_FILE=$NODE1_ORG3_CA -e CORE_PEER_ADDRESS=node1.organization3.gdzce.cn:11051 -e CORE_PEER_LOCALMSPID=Organization3MSP -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization3.gdzce.cn/users/Admin@organization3.gdzce.cn/msp cli  peer chaincode install -n ${CHAINCODE_NAME} -v ${CHAINCODE_VERSION} -l golang -p github.com/chaincode/${CHAINCODE_NAME}

echo "node1.organization2 安装链码"
docker exec -e CORE_PEER_TLS_ROOTCERT_FILE=$NODE1_ORG2_CA -e CORE_PEER_ADDRESS=node1.organization2.gdzce.cn:9051 -e CORE_PEER_LOCALMSPID=Organization2MSP -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization2.gdzce.cn/users/Admin@organization2.gdzce.cn/msp cli  peer chaincode install -n ${CHAINCODE_NAME} -v ${CHAINCODE_VERSION} -l golang -p github.com/chaincode/${CHAINCODE_NAME}

# 实例化链码
echo "实例化链码"
docker exec cli peer chaincode instantiate -o order.gdzce.cn:7050 -C channelname -n ${CHAINCODE_NAME} -v ${CHAINCODE_VERSION} -l golang -c '{"Args":["init"]}' -P 'AND("Organization3MSP.member","Organization1MSP.member","Organization2MSP.member")' --tls true --cafile $ORDERER_CA
sleep 10

# 调用链码查询
echo "调用链码"
docker exec -e CORE_PEER_TLS_ROOTCERT_FILE=$NODE1_ORG3_CA -e CORE_PEER_ADDRESS=node1.organization3.gdzce.cn:11051 -e CORE_PEER_LOCALMSPID=Organization3MSP -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization3.gdzce.cn/users/Admin@organization3.gdzce.cn/msp cli peer chaincode query -C channelname -n ${CHAINCODE_NAME} -c '{"Args":["getPapers",""]}'
docker exec -e CORE_PEER_TLS_ROOTCERT_FILE=$NODE1_ORG2_CA -e CORE_PEER_ADDRESS=node1.organization2.gdzce.cn:9051 -e CORE_PEER_LOCALMSPID=Organization2MSP -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization2.gdzce.cn/users/Admin@organization2.gdzce.cn/msp cli peer chaincode query -C channelname -n ${CHAINCODE_NAME} -c '{"Args":["getPapers",""]}'

echo -e "====== 链码 ${CHAINCODE_NAME} ${CHAINCODE_VERSION} 安装完成======"
