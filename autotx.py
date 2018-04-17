import os
import json
import sys
import random
import time


count = 1000
#account = "postquantum"
if len(sys.argv) > 1:
	count = int(sys.argv[1])

def formatJson(json):
	json = json.replace("\n", "")
	return json

def unlock(password):
	s = os.popen("hcashctl --wallet walletpassphrase " + password + " 0").read()

def getbalance():
	s = os.popen("hcashctl --wallet getbalance").read()
	balances = json.loads(formatJson(s))
	for balance in balances["balances"]:
		if balance["accountname"] == account:
			return balance["spendable"]

def tx(index, address, addrSize, amount, minconf):
	r = random.randint(0, addrSize - 1)
	#balance = getbalance()
	addr = address[r]
	
	amount = random.uniform(0, amount) + 0.1
	#if balance >= amount :
		#os.popen("hcashctl --wallet sendtoaddress " + addr + " " + str(amount)).read()
	p = os.system("hcashctl --wallet sendfrom " + account + " "   + addr + " " + str(amount) + " 0 " + str(minconf))
	
	print("[" + str(index) + "] send to " + addr + " amount: " + str(amount))

	return p


filename = "autotx.conf"

file = open(filename, "r")

s = formatJson(file.read())

conf = json.loads(s)

password = conf["password"]
address = conf["address"]
txAmount = conf["maxtxamount"]
minconf = conf["minconf"]
account = conf["account"]

addrSize = len(address)

balance = getbalance()

if txAmount <= 0:
	txAmount = balance / 1000

unlock(password)
print("unlocked")

print("Tx count: " + str(count))

#	for i in range(1,int(count)):
#		p = tx(i, address, addrSize, txAmount, minconf)
#		if p != 0:
#			break

if count == 0:
	count = -1

i = 0

startTime = time.time()

while i != count:
	p = tx(i, address, addrSize, txAmount, minconf)
	
#	if p != 0:
#		break

	i = i + 1

	#if balance == 0:
	#	print("No spendable money in your wallet.")
	#	break
endTime = time.time()

t = endTime - startTime
txPerSec = i / t
print("Time elapsed %.2fs" %  t) 
print("Send %.2f tx per second" % txPerSec) 
