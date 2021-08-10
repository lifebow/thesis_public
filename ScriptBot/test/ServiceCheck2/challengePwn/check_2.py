from pwn import *
from pprint import pprint

libc = ELF("/lib/x86_64-linux-gnu/libc.so.6")

def check_challenge(ip,port):
    r = remote(ip, port)
    r.recvuntil("system: ")
    system_addr = int(r.recvline().strip(), 16)
    libc.address = system_addr - libc.symbols["system"]

    ret = 0x40064e
    pop_rdi = 0x400963
    payload = "AAAA"
    r.sendline(payload)
    output =r.recvline().split("\n")[0]
    pprint(output)
    if output == "What's your name? Hello AAAA":
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


