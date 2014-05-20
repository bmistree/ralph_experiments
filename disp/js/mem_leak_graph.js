MEM_LEAK_TESTS_DIV_ID = 'mem_leak_tests';

/**
 @param {list} data_list --- Each element is a list containing
 integers and floats.  Indices are formatted according to 
 */
function draw_mem_leak_graph(data_list)
{
    // keys are number of operations performed, values are lists of
    // floats.  Each float is measured throughput of that condition.
    var mem_leak_graph_data_dict = {};
    // each element is the number of reads performed on a particular
    // condition.
    var num_reads_single_thread_list = [];
    for (var index in data_list)
    {
        var condition_data = new RunStats(data_list[index]);

        var condition_data_id =
            condition_data_id_generator(
                condition_data.num_reads,condition_data.num_threads);
                                        
        if (! (condition_data_id in mem_leak_graph_data_dict))
        {
            mem_leak_graph_data_dict[condition_data_id] = [];
            if (condition_data.num_threads == 1)
                num_reads_single_thread_list.push(condition_data.num_reads);
        }

        mem_leak_graph_data_dict[condition_data_id].push(condition_data);
    }

    num_reads_single_thread_list.sort(
        function (a,b)
        {
            if (a.num_reads < b.num_reads)
                return -1;
            if (a.num_reads > b.num_reads)
                return 1;
            return 0;
        }
    );
    var mem_leak_average_data_list = [];
    for (index in num_reads_single_thread_list)
    {
        var num_reads = num_reads_single_thread_list[index];
        var condition_data_id = condition_data_id_generator(num_reads,1);
        var avg_stats =
            stats_avg_throughput(mem_leak_graph_data_dict[condition_data_id]);
        mem_leak_average_data_list.push(avg_stats);
    }

    // grab two threads data and prepend it
    var two_threads_id =
        condition_data_id_generator(
            // note that we know that perform num_reads of smallest
            // with two threads.
            num_reads_single_thread_list[0],2);
    
    var avg_two_threads_reads =
        stats_avg_throughput(mem_leak_graph_data_dict[two_threads_id]);
    mem_leak_average_data_list.unshift(avg_two_threads_reads);

    plot(MEM_LEAK_TESTS_DIV_ID,mem_leak_average_data_list);
}

function condition_data_id_generator(num_reads_per_thread,num_threads)
{
    return num_reads_per_thread + "|" + num_threads;
}