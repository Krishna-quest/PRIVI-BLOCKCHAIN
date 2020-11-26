for i in {1..2}
do
    docker stop $(docker ps -q --filter name=dev-peer0.companies.com-CoinBalanceS$i)
done

for i in {1..28}
do
    docker stop $(docker ps -q --filter name=dev-peer0.companies.com-DataProtocolS$i)
done

for i in {1..12}
do
    docker stop $(docker ps -q --filter name=dev-peer0.companies.com-TraditionalLendingS$i)
done

for i in {1..29}
do
    docker stop $(docker ps -q --filter name=dev-peer0.companies.com-PRIVIcreditS$i)
done

for i in {1..115}
do
    docker stop $(docker ps -q --filter name=dev-peer0.companies.com-PodSwappingS$i)
done

for i in {1..44}
do
    docker stop $(docker ps -q --filter name=dev-peer0.companies.com-PodNFTS$i)
done