# Load Balancer

## Work flow based on code snippet

    ![This is a alt text.](./loadbalancer.png)


## Trade-offs:
    1. Using etcd as global variable map.
    2. Using etcd to store request references rather than worker object or time service object references.
    3. Health check routine not implemented due to restriction of sharing references.
    4. To the 3rd point limitation checking available worker while request to avoid panic in case of inavailability of worker node after kill.
    5. Using separate instance pool to hold registered worker object.
    6. Randomly allocating worker to serve a request. 
    7. Not secure.

## Improvements Needed:
    1. Including etcd as a member of load balancer object.
    2. Including instance pool as a member of load balancer to avoid using addition worker pool.
    3. Implementing separate health check goroutine which will be used to monitor health of worker or time-service which are registered in refernce to the point of improvement mentioned above.
    4. Implementing different schedular algorithms to schedule worker for request along with serving capacity in realtime. And using it on the basis of, at which level load balancer is to be used and in which manner(internal/external).
    5. Improving spawn controller event to register instance on the basis of config file feed.
    6. Change in method decleration syntax. Methods should be implemented using load balancer type.
    7. Enchancement in controller feature for security. 
        7.1 Token controller event: Use this token while spawning.
        7.2 Kill All: To remove all worker node which is obsolete, to avoid use of server resource in case of any network breach.


