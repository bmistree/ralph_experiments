NUM_READS_INDEX = 0;
NUM_THREADS_INDEX = 1;
READ_TYPE_INDEX = 2;
READS_PER_SECOND_INDEX = 3;

BAR_WIDTH = 40;
BAR_CATEGORY_SPACING_WIDTH = 10;
BAR_HEIGHT = 200;


function NonPerfStats(data_list)
{
    this.num_reads = data_list[NUM_READS_INDEX];
    this.num_threads = data_list[NUM_THREADS_INDEX];
    this.read_type = data_list[READ_TYPE_INDEX];
    this.reads_per_second = data_list[READS_PER_SECOND_INDEX];
}

NonPerfStats.prototype.clone = function()
{
    var data_list_init = [];
    data_list_init[NUM_READS_INDEX] = this.num_reads;
    data_list_init[NUM_THREADS_INDEX] = this.num_threads;
    data_list_init[READ_TYPE_INDEX] = this.read_type;
    data_list_init[READS_PER_SECOND_INDEX] = this.reads_per_second;
    return new NonPerfStats(data_list_init);
};

/**
 @param {list} stats_array --- Each element is a NonPerfStats object.
 
 @returns {NonPerfStats} --- Contains same num_reads, num_threads, and
 read_type as first object in stats_array.  reads_per_second is
 average value.
 */
function non_perf_stats_avg_throughput(stats_array)
{
    if (stats_array.length == 0)
        throw 'Cannot calulate average value of empty array';
    
    var avg = 0;
    for (var index in stats_array)
        avg += stats_array[index].reads_per_second;
    avg /= stats_array.length;


    var to_return = stats_array[0].clone();
    to_return.reads_per_second = avg;
    return to_return;
}
