# Multiple Locks with Key

Locking and Unlocking in golang using Mutual Exclusion ensure that only single goroutine is executing lines of code and other calling goroutines will have to wait until acquired lock is released by executing goroutine. 

When the lock is released by executing goroutine, other goroutine will be racing to acquire the lock, and block again when one goroutine has acquired the lock. 

While synchronization is achieved using simple Locking and Unlocking, it is inefficient for heavy operations that need advanced locking mechanism and could cause unnecessary bottleneck. 

Using multiple locks based on key, we can provide multiple locks to multiple goroutines with different keys. By doing so, multiple operations with different keys can execute codes simultaneously while also achieved synchronization for multiple operations using same key.
