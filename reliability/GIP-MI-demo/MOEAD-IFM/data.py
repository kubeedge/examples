import pandas as pd
import random

def create_simulation_data(edge_nodes,cloud_nodes, num_tasks):
    # 生成节点数据
    nodes = []
    for i in range(cloud_nodes):
        cpu_remaining = random.uniform(1000, 3000)  # CPU剩余量，假设范围在1000到8000之间，单位为MHz
        memory_remaining = random.uniform(1000, 8000)  # 内存剩余量，假设范围在2000到16000之间，单位为MB
        bandwidth = random.uniform(50, 200)  # 带宽，假设范围在50到200之间，单位为MB/s
        nodes.append(
            {'id': i, 'cpu_remaining': cpu_remaining, 'memory_remaining': memory_remaining, 'bandwidth': bandwidth,'type':"cloud"})

    for i in range(edge_nodes):
        cpu_remaining = random.uniform(500, 1500)  # CPU剩余量，假设范围在1000到8000之间，单位为MHz
        memory_remaining = random.uniform(500, 3000)  # 内存剩余量，假设范围在2000到16000之间，单位为MB
        bandwidth = random.uniform(50, 100)  # 带宽，假设范围在50到200之间，单位为MB/s
        nodes.append(
            {'id': i, 'cpu_remaining': cpu_remaining, 'memory_remaining': memory_remaining, 'bandwidth': bandwidth,
             'type': "edge"})

    # 生成任务数据
    tasks = []
    for i in range(num_tasks):
        cpu_demand = random.randint(1, 1000)  # CPU需求量，假设范围在100到1000之间，单位为MHz
        memory_demand = random.randint(50, 1500)  # 内存需求量，假设范围在200到1000之间，单位为MB
        container_cpu=0.1*random.randint(1, 10)
        container_mem=0.1*random.randint(1, 10)
        priority = random.choice(["High", "Medium", "Low"])  # 任务优先级
        data_size = random.uniform(1, 100)  # 数据大小，假设范围在1到100之间，单位为MB
        exec = random.uniform(0.1, 2)  # 任务处理时间，假设范围在0.5.0到2.0之间，单位为秒
        tasks.append({'id': i, 'priority': priority,'cpu_demand': cpu_demand, 'memory_demand': memory_demand,
                       'data_size': data_size, 'exec': exec,'container_cpu':container_cpu,'container_mem':container_mem})

    # #生成节点数据
    # nodes = []
    # for i in range(num_nodes):
    #     cpu_remaining = random.uniform(500, 2000)  # CPU剩余量，假设范围在1000到8000之间，单位为MHz
    #     memory_remaining = random.uniform(500, 5000)  # 内存剩余量，假设范围在2000到16000之间，单位为MB
    #     bandwidth = random.uniform(50, 100)  # 带宽，假设范围在50到200之间，单位为MB/s
    #     nodes.append(
    #         {'id': i, 'cpu_remaining': cpu_remaining, 'memory_remaining': memory_remaining, 'bandwidth': bandwidth})
    #
    # # 生成任务数据
    # tasks = []
    # for i in range(num_tasks):
    #     cpu_demand = random.randint(100, 1000)  # CPU需求量，假设范围在100到1000之间，单位为MHz
    #     memory_demand = random.randint(50, 2000)  # 内存需求量，假设范围在200到1000之间，单位为MB
    #     container_cpu=random.randint(1, 100)
    #     container_mem=random.randint(1, 100)
    #     priority = random.choice(["High", "Medium", "Low"])  # 任务优先级
    #     data_size = random.uniform(50, 1000)  # 数据大小，假设范围在100到1000之间，单位为MB
    #     exec = random.uniform(0.1, 2)  # 任务处理时间，假设范围在0.5.0到2.0之间，单位为秒
    #     tasks.append({'id': i, 'priority': priority,'cpu_demand': cpu_demand, 'memory_demand': memory_demand,
    #                    'data_size': data_size, 'exec': exec,'container_cpu':container_cpu,'container_mem':container_mem})


    # 转换为DataFrame
    nodes_df = pd.DataFrame(nodes)
    tasks_df = pd.DataFrame(tasks)

    # 保存到CSV文件
    nodes_df.to_csv('nodes0401.csv', index=False)
    tasks_df.to_csv('tasks0401.csv', index=False)

if __name__ == '__main__':
    # 生成数据
    create_simulation_data(200,100, 30)
