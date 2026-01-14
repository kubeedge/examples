import os
import random

import numpy as np
import pandas as pd

tasks = None
nodes = None


class UnloadAssignment:
    def __init__(self, task, node):
        self.task = task
        self.node = node
        self.fitness = [0, 0, 0, 0]  # 响应时间, 资源适应度, 资源释放收益, 任务依赖度

    def calculate_fitness(self):
        # 目标1: 最小化任务响应时间
        total_delay = self.cal_delay()

        # 目标2: 最大化资源适应度（取负转为最小化）
        node_fitness = self.cal_nodeFit()

        # 目标3: 最大化资源释放收益（取负转为最小化）
        resource_release = self.cal_resourceRelease()

        # 目标4: 最小化任务依赖度
        task_dependency = self.cal_taskDependency()

        self.fitness = [total_delay, -node_fitness, -resource_release, task_dependency]

    def cal_delay(self):
        node = nodes[self.node]
        task = tasks[self.task]
        transmission_delay = task.data_size / node.bandwidth
        return transmission_delay + task.exec

    def cal_nodeFit(self):
        cpu_fit = (nodes[self.node].cpu_remaining - tasks[self.task].cpu_demand) / nodes[self.node].cpu_remaining
        mem_fit = (nodes[self.node].memory_remaining - tasks[self.task].memory_demand) / nodes[
            self.node].memory_remaining
        return 0.6 * cpu_fit + 0.4 * mem_fit  # 保持原负值

    def cal_resourceRelease(self):
        task = tasks[self.task]
        return 0.5 * task.container_cpu + 0.5 * task.container_mem  # 资源释放收益

    def cal_taskDependency(self):
        if nodes[self.node].type == 'cloud':
            if tasks[self.task].priority == "High":
                return 10  # 需要最小化的值
            elif tasks[self.task].priority == "Medium":
                return 1
            else:
                return 0
        if nodes[self.node].type == 'edge':
            if tasks[self.task].priority == "High":
                return 1  # 需要最小化的值
            elif tasks[self.task].priority == "Medium":
                return 0
            else:
                return 0


class Node:
    def __init__(self, id, cpu_remaining, memory_remaining, bandwidth, type):
        self.id = id
        self.cpu_remaining = cpu_remaining
        self.memory_remaining = memory_remaining
        self.bandwidth = bandwidth
        self.type = type


class Task:
    def __init__(self, id, priority, cpu_demand, memory_demand, data_size, exec, container_cpu, container_mem):
        self.id = id
        self.priority = priority
        self.cpu_demand = cpu_demand
        self.memory_demand = memory_demand
        self.data_size = data_size
        self.exec = exec
        self.container_cpu = container_cpu
        self.container_mem = container_mem


class Para:
    def __init__(self):
        self.N_elite = 2
        self.Pop_size = 50  # 增大种群规模适应四维目标
        self.gene_size = 500
        self.pc_max = 0.8
        self.pm_max = 0.1
        self.pc_min = 0.7
        self.pm_min = 0.01
        self.pc = 0.8
        self.pm = 0.3
        self.T = 5
        self.p_GS = 0.5
        self.p_LS = 0.3
        self.p_RS = 0.2
        self.Best_JS = None
        self.Best_Cmax = 1e+20
        self.C_end = []
        self.Pop = []
        self._lambda = FourD_VGM(5)  # 四维权重生成
        self._z = []
        # IFM
        self.history_pop = []


def FourD_VGM(H):
    delta = 1 / H
    weights = []
    for w1 in np.arange(0, 1 + delta, delta):
        for w2 in np.arange(0, 1 - w1 + delta, delta):
            for w3 in np.arange(0, 1 - w1 - w2 + delta, delta):
                w4 = 1 - w1 - w2 - w3
                if w4 >= -1e-9:  # 处理浮点误差
                    weights.append([w1, w2, w3, w4])
    return weights


def Tri_Dominate(Pop1, Pop2):
    return all(x <= y for x, y in zip(Pop1.fitness, Pop2.fitness)) and any(
        x < y for x, y in zip(Pop1.fitness, Pop2.fitness))


def Neighbor(lambd, T):
    B = []
    for i in range(len(lambd)):
        temp = []
        for j in range(len(lambd)):
            distance = np.linalg.norm(np.array(lambd[i]) - np.array(lambd[j]))
            temp.append(distance)
        res = sorted(range(len(temp)), key=lambda k: temp[k])[:T]
        B.append(res)
    return B


def Tchebycheff(x, z, lambd):
    Gte = [abs(xi - zi) * wi for xi, zi, wi in zip(x.fitness, z, lambd)]
    return max(Gte)


# 修改参考点计算
def find_worst_case(all_solutions):
    max_vals = [
        max(sol.fitness[0] for sol in all_solutions),
        max(sol.fitness[1] for sol in all_solutions),
        max(sol.fitness[2] for sol in all_solutions),
        max(sol.fitness[3] for sol in all_solutions)
    ]
    return [v + 1 for v in max_vals]


def check_constraint(task_index, node_index):
    task = tasks[task_index]
    node = nodes[node_index]
    return node.cpu_remaining >= task.cpu_demand and node.memory_remaining >= task.memory_demand


def read_simulation_data():
    base_dir = os.path.dirname(os.path.abspath(__file__))
    nodes_df = pd.read_csv(os.path.join(base_dir, '仿真数据', 'nodes0401.csv'))
    tasks_df = pd.read_csv(os.path.join(base_dir, '仿真数据', 'tasks0401.csv'))
    nodes = [Node(row['id'], row['cpu_remaining'], row['memory_remaining'], row['bandwidth'], row['type']) for
             index, row in nodes_df.iterrows()]
    tasks = [Task(row['id'], row['priority'], row['cpu_demand'], row['memory_demand'], row['data_size'], row['exec'],
                  row['container_cpu'], row['container_mem']) for index, row in tasks_df.iterrows()]
    return nodes, tasks


def get_solutions(nodes, tasks):
    solutions = []
    task_num = len(tasks)
    node_num = len(nodes)
    for i in range(task_num):
        for j in range(node_num):
            if check_constraint(i, j):
                assignment = UnloadAssignment(i, j)
                assignment.calculate_fitness()  # 计算适应度
                solutions.append(assignment)
    return solutions


def true_pareto_frontier(solutions):
    pareto_frontier = []
    for solution in solutions:
        is_pareto = True
        for other_solution in solutions:
            if solution != other_solution:
                if all(y <= x for x, y in zip(solution.fitness, other_solution.fitness)) and any(
                        y < x for x, y in zip(solution.fitness, other_solution.fitness)):
                    is_pareto = False
                    break
        if is_pareto:
            pareto_frontier.append(solution)
    return pareto_frontier


def crossover(parent1, parent2):
    new_pop1 = parent1
    new_pop2 = parent2
    # 从父代中随机选择任务和节点
    task1, task2 = parent1.task, parent2.task
    node1, node2 = parent1.node, parent2.node

    child1 = UnloadAssignment(task1, node2)
    child2 = UnloadAssignment(task2, node1)

    # 添加约束检查，确保子代满足约束条件
    if check_constraint(child1.task, child1.node):
        new_pop1 = child1
    if check_constraint(child2.task, child2.node):
        new_pop2 = child2
    return new_pop1, new_pop2


def mutation(solution, num_tasks, num_nodes):
    task = solution.task
    node = solution.node
    while True:
        if random.random() < 0.5:
            task = random.randint(0, num_tasks - 1)
        else:
            node = random.randint(0, num_nodes - 1)
        # 添加约束检查，确保目标节点的CPU和内存剩余量大于要卸载的任务的资源需求量
        if check_constraint(task, node):
            return UnloadAssignment(task, node)


def operator(chs1, chs2, pc, pm, num_tasks, num_nodes):
    p1, p2 = chs1, chs2
    if random.random() < pc:  # 是否进行交叉
        p1, p2 = crossover(chs1, chs2)

    if random.random() < pm:  # 是否进行变异
        p1 = mutation(p1, num_tasks, num_nodes)

    if random.random() < pm:
        p2 = mutation(p2, num_tasks, num_nodes)

    # p1 = mutation(p1, num_tasks, num_nodes)
    # p2 = mutation(p2, num_tasks, num_nodes)

    return p1, p2


def MOEAD_main(para):
    num_nodes = len(nodes)
    num_tasks = len(tasks)

    para.Pop_size = len(para._lambda)
    para.Pop = []

    # 初始化种群
    for _ in range(para.Pop_size):
        assignment = None
        while (True):
            task_id = random.randint(0, num_tasks - 1)
            node_id = random.randint(0, num_nodes - 1)
            if check_constraint(task_id, node_id):
                assignment = UnloadAssignment(task_id, node_id)
                break
        assignment.calculate_fitness()  # 计算适应度
        para.Pop.append(assignment)
        if para._z == []:
            para._z = assignment.fitness
        else:
            for i in range(len(assignment.fitness)):
                if assignment.fitness[i] < para._z[i]:
                    para._z[i] = assignment.fitness[i]

    B = Neighbor(para._lambda, para.T)  # 计算邻居
    EP = []  # 存储非支配解
    convergence_speed = []
    for gi in range(para.gene_size):
        # ==== IFM新增：保存历史种群 ====
        if len(para.history_pop) >= 1:  # 只保留上一代
            para.history_pop.pop(0)
        para.history_pop.append([ind for ind in para.Pop])  # 深拷贝当前种群
        for i in range(len(para.Pop)):
            # 从B(i)中随机选择两个索引
            j = random.randint(0, para.T - 1)
            k = random.randint(0, para.T - 1)
            # 生成新解
            pop1, pop2 = operator(para.Pop[B[i][j]], para.Pop[B[i][k]], para.pc, para.pm, num_tasks, num_nodes)
            pop1.calculate_fitness()  # 计算适应度
            pop2.calculate_fitness()  # 计算适应度

            if Tri_Dominate(pop1, pop2):
                y = pop1
            else:
                y = pop2

            # ==== IFM模型（适配本场景） ====
            if gi >= 1:
                prev_pop = para.history_pop[-1]
                prev_ind = prev_pop[i]

                # 计算切比雪夫适应度
                lambda_i = para._lambda[i]  # 当前子问题的权重
                T_y = Tchebycheff(y, para._z, lambda_i)
                T_prev = Tchebycheff(prev_ind, para._z, lambda_i)

                if T_y > T_prev:
                    # 执行交叉操作
                    new_y1, new_y2 = crossover(y, prev_ind)

                    # 计算新个体适应度
                    new_y1.calculate_fitness()
                    new_y2.calculate_fitness()

                    # 选择更好的个体
                    T_new1 = Tchebycheff(new_y1, para._z, lambda_i)
                    T_new2 = Tchebycheff(new_y2, para._z, lambda_i)
                    y_new = new_y1 if T_new1 < T_new2 else new_y2

                    # # 约束检查
                    if check_constraint(y_new.task, y_new.node):
                        y_new.calculate_fitness()
                        # T_y_new = Tchebycheff(y_new, para._z, lambda_i)
                        # y = y_new if T_y_new < T_y else y
                        if Tri_Dominate(y_new, y):
                            y=y_new

            # 更新参考点
            for zi in range(len(para._z)):
                if para._z[zi] > y.fitness[zi]:
                    para._z[zi] = y.fitness[zi]
            # 更新邻居解
            for bi in range(len(B[i])):
                Ta = Tchebycheff(para.Pop[B[i][bi]], para._z, para._lambda[B[i][bi]])
                Tb = Tchebycheff(y, para._z, para._lambda[B[i][bi]])
                if Tb < Ta:
                    para.Pop[B[i][bi]] = y
            # 更新EP
            if EP == []:
                EP.append(y)
            else:
                dominateY = False  # 是否有支配Y的解
                equalFlg = False
                _remove = []  # 从EP中删除所有被Y支配的解
                for ei in range(len(EP)):
                    if Tri_Dominate(y, EP[ei]):
                        _remove.append(EP[ei])
                    elif Tri_Dominate(EP[ei], y):
                        dominateY = True
                        break
                    elif y.node == EP[ei].node and y.task == EP[ei].task:
                        equalFlg = True
                        break
                # 如果EP中没有向量支配Y，则将Y添加到EP中
                if not dominateY and not equalFlg:
                    EP.append(y)
                    for j in range(len(_remove)):
                        EP.remove(_remove[j])
        convergence_speed.append(len(EP))
    return EP, convergence_speed


def calculate_IGD(P, PF):
    total_distance = 0
    for p in PF:
        min_dist = min(np.sqrt(sum((p.fitness[i] - q.fitness[i]) ** 2 for i in range(len(p.fitness)))) for q in P)
        total_distance += min_dist
    return total_distance / len(PF)


def calculate_hypervolume(P, reference_point):
    P = [list(sol.fitness) for sol in P]
    try:
        from deap.tools.indicator import hv
    except Exception as e:
        raise RuntimeError("deap is required for hypervolume calculation. Install it via `pip install deap`.") from e
    return hv.hypervolume(P, reference_point)


def main():
    global nodes, tasks
    nodes, tasks = read_simulation_data()

    all_solutions = get_solutions(nodes, tasks)
    true_pf = true_pareto_frontier(all_solutions)
    print("true_pareto_frontier", len(true_pf))

    para = Para()
    EP, _convergence_speed = MOEAD_main(para)

    reference_point = find_worst_case(all_solutions)
    try:
        hv_value = calculate_hypervolume(EP, reference_point)
        print(f"Hypervolume: {hv_value:.4f}")
    except RuntimeError:
        pass

    igd_value = calculate_IGD(EP, true_pf)
    print(f"IGD: {igd_value:.4f}")
    print(f"Pareto solutions: {len(EP)}")


if __name__ == '__main__':
    main()

