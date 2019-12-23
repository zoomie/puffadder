Key value store where the key is a persons name and the value is an integer amount. 

Current Todo:
- Have the data persisted on disk (like the D in ACID databases). 
- Be able to configure at least two different algorithms to hold the Index (I'm currently calling the offset of the bytes.)
- Implement a simple transaction so that you could either move money from one account and into the other in an atomic step, or it doesn't happen at all. 