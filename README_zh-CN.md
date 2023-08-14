# pterergate-dtf

Pterergate-dtf (Pterergate Distributed Task Framework) 是一款高性能的分布式任务调度框架，支持在由数万个节点组成的集群中并行调度数千个正在运行的任务。

pterergate-dtf支持把任务拆分成若干个更小粒度、耗时更短的子任务，并将拆分出的子任务调度到集群的节点上执行。一般来讲，可以将一个任务拆分成数十到数千个可以并行执行的子任务，以便充分利用集群的规模优势，提高任务的执行效率，减少执行耗时。
