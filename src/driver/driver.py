import json
import random
import math
import bridge

def generate_interval_n_size_sequence_by_json(filename: str) -> (dict):
    """Generate interval and size sequence for the users inside the cases
    
    Generate interval and size sequence for the users inside the cases, this function directly accept the json file, and do parsing
    to fetch key-value pairs in json, and generate python dictionary for loading these data.
    This function is going to read the dictionary while doing the interval and size sequence generation for each user in each case. 
    
    Args:
        filename: the file name of the json file, which must be obey to the sample json format
    
    Returns:
        A system generated dictionary which include the system generated interval sequence and size sequence for each case
        
        {Case_id: {
            User_id: {
                'interval_seq': [interval1, interval2, interval3, ...],
                'size_seq': [size1, size2, size3, ... ],
                'user_ip': '127.0.0.1'
                }}}
        
        Returned dictionary key can be integer when it is case_id and user_id, or string when it is the key for fetching 
        interval sequence and time sequence and the value of the key can be dictionary or list depends on the layer to be obtained
    """
    
    # Opening JSON file
    f = open(filename)

    # returns JSON object as a dictionary
    data = json.load(f)
    interval_n_size_sequence = {} # user_behavior_dict['userid'] = behavior dict
    
    for case in data['cases']:
        case_id = case['case_id']
        timeline = case['timeline']
        server_mode = case['server_mode']
        server_ip = case['server_ip']
        
        tmp_intv_n_size_seq_dict = {}
        for user in case['users']:
            user_id = user['userid']
            user_ip = user['userip']
            behavior = {
                'interval_seq_pattern': user['interval_seq_pattern'],
                'min_interval': user['min_interval'],
                'max_interval': user['max_interval'],
                'custom_interval_seq': user['custom_interval_seq'],
                'protocol': user['protocol'],
                'size_seq_pattern': user['size_seq_pattern'],
                'min_size': user['min_size'],
                'max_size': user['max_size'],
                'custom_size_seq': user['custom_size_seq'],
                'disconnection_enable': user['disconnection_enable'],
                'disconnection_prob': user['disconnection_prob'],
                'disconnection_mode': user['disconnection_mode'],
                'max_prob': user['max_prob'],
                'min_prob': user['min_prob'],
                "disconnect_time_pattern": user['disconnect_time_pattern'],
                "disconnect_time_seq": user['disconnect_time_seq'],
                "disconnect_min_time": user['disconnect_min_time'],
                "disconnect_max_time": user['disconnect_max_time'],   
            }
            
            tmp_intv_n_size_seq_dict[user_id] = generate_interval_n_size_sequence(timeline,behavior)
            tmp_intv_n_size_seq_dict[user_id]['user_ip'] = user_ip
            
        interval_n_size_sequence[case_id] = tmp_intv_n_size_seq_dict
        interval_n_size_sequence[case_id]['server_mode'] = server_mode
        interval_n_size_sequence[case_id]['server_ip'] = server_ip

    # Closing file
    f.close()
    return interval_n_size_sequence


def generate_interval_n_size_sequence(timeline: int,behavior: dict) -> dict:
    """Generate interval and size sequence
    
    Generate interval sequence and size sequence depends on the user behavior dictionary
    
    Args:
        timeline: the total time to be consumed for the test case
        behavior: the user behavior defined in the json file
    
    Returns:
        A system generated dictionary which include the system generated interval sequence and size sequence
        
        {'interval_seq': [interval1, interval2, interval3, ...],
            'size_seq': [size1, size2, size3, ... ]}
        
        Returned key is always string, and the value of the key is always list.
        The returned interval sequence list and size sequence list is the same size
    """
    
    tmp_intv_n_size_dict = {}
    
    # disconnection related params
    
    disconnect_behavior_dict = {}
    
    disconnect_enable = behavior['disconnection_enable']
    disconnection_prob = behavior['disconnection_prob']
    disconnection_mode = behavior['disconnection_mode']
    max_prob = behavior['max_prob']
    min_prob = behavior['min_prob']
    
    ## TODO: To be implemented

    disconnect_time_pattern = behavior['disconnect_time_pattern']
    disconnect_time_seq = behavior['disconnect_time_seq']
    disconnect_min_time = behavior['disconnect_min_time']
    disconnect_max_time = behavior['disconnect_max_time']

    
    
    # interval related params
    interval_seq_pattern = behavior['interval_seq_pattern']
    min_i = behavior['min_interval']
    max_i = behavior['max_interval']
    custom_interval_seq = behavior['custom_interval_seq']

    # size related params
    size_seq_pattern = behavior['size_seq_pattern']
    min_s = behavior['min_size']
    max_s = behavior['max_size']
    custom_size_seq = behavior['custom_size_seq']

    # protocol related params
    protocol = behavior['protocol']
    
    
    # generate interval sequence pattern
    custom_interval_sum = sum(custom_interval_seq)
    tmp_intv_n_size_dict["interval_seq"] = generate_interval_sequence(interval_seq_pattern, timeline, custom_interval_sum, 
                                                                          custom_interval_seq , min_i, max_i)

    # generate size sequence pattern
    final_interval_seq_len = len(tmp_intv_n_size_dict['interval_seq'])
    
    size_seq_ret_dict = generate_size_sequence(size_seq_pattern, timeline, custom_interval_sum,custom_size_seq, 
                                               final_interval_seq_len, min_s, max_s,
                                               disconnect_enable, disconnection_prob, 
                                               disconnection_mode, max_prob, min_prob)
    tmp_intv_n_size_dict['size_seq'] = size_seq_ret_dict['size_seq']

    # generate disconnection time sequence pattern
    required_size = size_seq_ret_dict['disconnect_cnt']
    tmp_intv_n_size_dict['disconnect_time_seq'] = generate_disconnect_time_sequence(disconnect_time_pattern, disconnect_time_seq, 
                                                                                    disconnect_min_time, disconnect_max_time, 
                                                                                    required_size)
        
    return tmp_intv_n_size_dict        


def generate_interval_sequence(pattern: str, timeline: int, interval_sum: int, custom_interval_list: list,
                               min_i: int, max_i: int) -> list:
    """Generate a interval sequence
    
    Generate a interval sequence depends on different pattern type
    
    Args:
        pattern: a string that indicate the pattern for size sequence generation
        timeline: the total time to be consumed for the test case
        interval_sum: the sum of elements of user customized interval sequence
        custom_interval_list: the user customized interval sequence which can be fetched from json file
        min_s: the min size of the frame -> defined in json file
        max_s: the max size of the frame -> defined in json file
    
    Returns:
        A system generated size sequence list based on indicated pattern
        example:
        
        [interval1, interval2, interval3, ...]
        
        Returned element inside the list is always integer, the integer is to indicate the interval to send frame by "ms"
        
    Raises:
        ValueError: An error occurred when giving wrong pattern string
    """
    
    # generate interval sequence pattern
    interval_seq = []
    
    if pattern == 'custom':
        total_repeat = math.ceil(timeline / interval_sum) # take ceil of float as total replication
        interval_seq = custom_interval_list * total_repeat

    elif pattern == 'random':
        tmp_total_intv = 0
        while tmp_total_intv < timeline:
            tmp_intv = random.randrange(min_i, max_i + 1)
            tmp_total_intv += tmp_intv
            interval_seq.append(tmp_intv)
            
    elif pattern == 'equal_diff_increment':
        a1 = min_i
        an = max_i
        Sn = timeline
        
        num_of_interval = math.ceil((Sn * 2)/(a1+an)) #Sn=n(a1+an)/2 

        increment = (max_i - min_i) / num_of_interval
        cumulative_increment = 0
        current_interval = min_i
        for cnt in range(0, num_of_interval):

            if cumulative_increment >= 1: # 1 ms -> min time interval to increase
                current_interval += int(cumulative_increment)
                cumulative_increment = 0
                
            cumulative_increment += increment

            interval_seq.append(current_interval)
            
    else:
        raise ValueError('Invalid size pattern mode !!!')
    
    return interval_seq


def generate_size_sequence(pattern: str, timeline: int, interval_sum: int, custom_size_list: list, 
                           interval_seq_len: int, min_s: int, max_s: int, disconnect_enable: bool, disconnection_prob: float,
                           disconnection_mode: str, max_prob: float, min_prob: float) -> dict:
    """Generate a size sequence
    
    Generate a size sequence depends on different pattern type, this function must be used after 
    the interval sequence has been generated by system
    
    Args:
        pattern: a string that indicate the pattern for size sequence generation
        timeline: the total time to be consumed for the test case
        interval_sum: the sum of elements of user customized interval sequence
        custom_size_list: the user customized size sequence which can be fetched from json file
        interval_seq_len: the length of system-generated interval sequence
        min_s: the min size of the frame -> defined in json file
        max_s: the max size of the frame -> defined in json file
    
    Returns:
        A system generated size sequence list based on indicated pattern
        example:
        
        [size1, size2, size3, ...]
        
        Returned element inside the list is always integer, the integer is to indicate the bytes of the frame
        
    Raises:
        ValueError: An error occurred when giving wrong pattern string
    """
    # define return dict
    ret_dict = {}
    ret_dict['size_seq'] = []
    ret_dict['disconnect_cnt'] = 0
    
    # generate size sequence pattern
    size_seq = []
#     print("interval {}".format(interval_seq_len))
    
    if pattern == 'custom':
        if disconnect_enable:
            circular_cnt = 0
            for cnt in range(0, interval_seq_len):
                if circular_cnt >= len(custom_size_list):
                    circular_cnt = 0
                tmp_size = custom_size_list[circular_cnt]
                size_seq.append(tmp_size) 
                circular_cnt += 1
                
                 # Set disconnect prob.
                trigger_discon = random_number_generator_with_prob(disconnection_mode, disconnection_prob, min_prob, max_prob)
                if trigger_discon == 1: # disconnection triggered
                    size_seq.append(-1)
                    ret_dict['disconnect_cnt'] += 1
                
        else:
            size_seq = custom_size_list * interval_seq_len
            size_seq = size_seq[: interval_seq_len]
        
    elif pattern == 'random':
        if disconnect_enable: # If disconnection enabled
            for cnt in range(0, interval_seq_len):
                
                tmp_size = random.randrange(min_s, max_s + 1)
                size_seq.append(tmp_size)  
                
                # Set disconnect prob.
                trigger_discon = random_number_generator_with_prob(disconnection_mode, disconnection_prob, min_prob, max_prob)    
                
                if trigger_discon == 1: # disconnection triggered
                    size_seq.append(-1)
                    ret_dict['disconnect_cnt'] += 1

        else:
            for cnt in range(0, interval_seq_len):
                tmp_size = random.randrange(min_s, max_s + 1)
                size_seq.append(tmp_size) 
            
    elif pattern == 'equal_diff_increment':
        increment = (max_s - min_s) / interval_seq_len
        cumulative_increment = 0
        current_size = min_s
        
        if disconnect_enable:
            for cnt in range(0, interval_seq_len):

                if cumulative_increment < 1:
                    cumulative_increment += increment

                else:
                    current_size += int(cumulative_increment)
                    cumulative_increment = 0
                    cumulative_increment += increment
                size_seq.append(current_size)
                
                # Set disconnect prob.
                trigger_discon = random_number_generator_with_prob(disconnection_mode, disconnection_prob, min_prob, max_prob)
                
                if trigger_discon == 1: # disconnection triggered
                    size_seq.append(-1)
                    ret_dict['disconnect_cnt'] += 1
                
        else:
            for cnt in range(0, interval_seq_len):

                if cumulative_increment < 1:
                    cumulative_increment += increment

                else:
                    current_size += int(cumulative_increment)
                    cumulative_increment = 0
                    cumulative_increment += increment

                size_seq.append(current_size)
    else:
        raise ValueError('Invalid size pattern mode !!!')
    
    ret_dict['size_seq'] = size_seq
    
    return ret_dict


def generate_disconnect_time_sequence(pattern: str, disconnect_time_seq: list, 
                                      min_t: int, max_t: int, required_size: int) -> list:
    disconnect_time_seq = []
    
    ## TODO: To be implemented
    if pattern == 'custom':
        total_repeat = math.ceil(required_size / len(disconnect_time_seq))
        disconnect_time_seq = (disconnect_time_seq * total_repeat)[:required_size]

    elif pattern == 'random':
        for i in range(0,required_size):
            disconnect_time_seq.append(random.randrange(min_t, max_t+1))
        
    elif pattern == 'equal_diff_increment':

        increment = (max_t - min_t) / required_size
        cumulative_increment = 0
        current_interval = min_t
        for cnt in range(0, required_size):

            if cumulative_increment >= 1: # 1 ms -> min time interval to increase
                current_interval += int(cumulative_increment)
                cumulative_increment = 0
                
            cumulative_increment += increment

            disconnect_time_seq.append(current_interval)
            
    else:
        raise ValueError('Invalid disconnect time pattern mode !!!')
    
    return disconnect_time_seq
    

def random_number_generator_with_prob(disconnection_mode: str, disconnection_prob: float, min_prob: float, max_prob: float)-> int:
    if disconnection_mode == 'random':
        min_prob = int(min_prob*100)
        max_prob = int(max_prob*100)
        random_prob = random.randrange(max_prob, max_prob+1) # 1% to 10% disconnect prob.
        trigger_discon = random.randrange(0, random_prob)
        return trigger_discon
    
    elif disconnection_mode == 'fixed':
        prob = int(disconnection_prob * 100)
        trigger_discon = random.randrange(0, prob+1)
        return trigger_discon
    
    else: 
        raise ValueError('Invalid disconnection mode !!!')
    
# Entry point here
if __name__ == '__main__':
    filename = "sample_json.json" # plug your json file name here
    ub_dict = generate_interval_n_size_sequence_by_json(filename)
    bridge.accept_data_from_driver(ub_dict)
