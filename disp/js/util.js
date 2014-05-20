NUM_READS_INDEX = 0;
NUM_THREADS_INDEX = 1;
READ_TYPE_INDEX = 2;
READS_PER_SECOND_INDEX = 3;

CPUS_UTILIZED_INDEX = 4;
CONTEXT_SWITCHES_INDEX = 5;
CPU_MIGRATIONS_INDEX = 6;
PAGE_FAULTS_INDEX = 7;
CYCLES_INDEX = 8;
STALLED_CYCLES_FRONTEND_INDEX = 9;
FRONTEND_CYCLES_IDLE_INDEX = 10;
STALLED_CYCLES_BACKEND_INDEX = 11;
BACKEND_CYCLES_IDLE_INDEX = 12;
INSTRUCTIONS_INDEX = 13;
BRANCHES_INDEX = 14;
BRANCH_MISSES_INDEX = 15;

BAR_WIDTH = 40;
BAR_CATEGORY_SPACING_WIDTH = 10;
BAR_HEIGHT = 200;

/**
 @param {string} div_id_to_plot_on
 @param {list} data_list --- Each element is a RunStats obj
 */
function plot(div_id_to_plot_on,data_list)
{
    var num_conditions = data_list.length;
    var width =
        (BAR_WIDTH + BAR_CATEGORY_SPACING_WIDTH)*num_conditions;
    
    var bar_chart = d3.select('#' + div_id_to_plot_on).
        append('svg:svg').
        attr('width', width).
        attr('height', BAR_HEIGHT);

    var x_rect_positions =
        d3.scale.linear().domain([0, num_conditions]).range([0, width]);
    var y_heights =
        d3.scale.linear().domain(
            [0, d3.max(data_list,
                       function(datum) { return datum.reads_per_second; })]).
        rangeRound([0, BAR_HEIGHT]);
    
    // FIXME: don't understand semantics of selectAll
    bar_chart.selectAll('rect').
        data(data_list).
        enter().
        append("svg:rect").
        attr("x",
             function(datum, index)
             {
                 return x_rect_positions(index);
             }).
        attr("y",
             function(datum)
             {
                 return BAR_HEIGHT - y_heights(datum.reads_per_second);
             }).
        attr("height",
             function(datum)
             {
                 return y_heights(datum.reads_per_second);
             }).
        attr("width", BAR_WIDTH).
        attr("fill", "steelblue");
}

function RunStats(data_list)
{
    this.num_reads = data_list[NUM_READS_INDEX];
    this.num_threads = data_list[NUM_THREADS_INDEX];
    this.read_type = data_list[READ_TYPE_INDEX];
    this.reads_per_second = data_list[READS_PER_SECOND_INDEX];

    this.cpus_utilized = null;
    this.context_switches = null;
    this.cpu_migrations = null;
    this.page_faults = null;
    this.cycles = null;
    this.stalled_cycles_frontend = null;
    this.frontend_cycles_idle = null;
    this.stalled_cycles_backend = null;
    this.backend_cycles_idle = null;
    this.instructions = null;
    this.branches = null;
    this.branch_misses = null;
    
    if (data_list.length > 4)
    {
        this.cpus_utilized = data_list[CPUS_UTILIZED_INDEX];
        this.context_switches = data_list[CONTEXT_SWITCHES_INDEX];
        this.cpu_migrations = data_list[CPU_MIGRATIONS_INDEX];
        this.page_faults = data_list[PAGE_FAULTS_INDEX];
        this.cycles = data_list[CYCLES_INDEX];
        this.stalled_cycles_frontend = data_list[STALLED_CYCLES_FRONTEND_INDEX];
        this.frontend_cycles_idle = data_list[FRONTEND_CYCLES_IDLE_INDEX];
        this.stalled_cycles_backend = data_list[STALLED_CYCLES_BACKEND_INDEX];
        this.backend_cycles_idle = data_list[BACKEND_CYCLES_IDLE_INDEX];
        this.instructions = data_list[INSTRUCTIONS_INDEX];
        this.branches = data_list[BRANCHES_INDEX];
        this.branch_misses = data_list[BRANCH_MISSES_INDEX];
    }
}

RunStats.prototype.clone = function()
{
    var data_list_init = [];
    data_list_init[NUM_READS_INDEX] = this.num_reads;
    data_list_init[NUM_THREADS_INDEX] = this.num_threads;
    data_list_init[READ_TYPE_INDEX] = this.read_type;
    data_list_init[READS_PER_SECOND_INDEX] = this.reads_per_second;

    data_list_init[CPUS_UTILIZED_INDEX] = this.cpus_utilized;
    data_list_init[CONTEXT_SWITCHES_INDEX] = this.context_switches;
    data_list_init[CPU_MIGRATIONS_INDEX] = this.cpu_migrations;
    data_list_init[PAGE_FAULTS_INDEX] = this.page_faults;
    data_list_init[CYCLES_INDEX] = this.cycles;
    data_list_init[STALLED_CYCLES_FRONTEND_INDEX] = this.stalled_cycles_frontend;
    data_list_init[FRONTEND_CYCLES_IDLE_INDEX] = this.frontend_cycles_idle;
    data_list_init[STALLED_CYCLES_BACKEND_INDEX] = this.stalled_cycles_backend;
    data_list_init[BACKEND_CYCLES_IDLE_INDEX] = this.backend_cycles_idle;
    data_list_init[INSTRUCTIONS_INDEX] = this.instructions;
    data_list_init[BRANCHES_INDEX] = this.branches;
    data_list_init[BRANCH_MISSES_INDEX] = this.branch_misses;
    
    return new RunStats(data_list_init);
};

/**
 @param {list} stats_array --- Each element is a NonPerfStats object.
 
 @returns {RunStats} --- Contains same num_reads, num_threads, and
 read_type as first object in stats_array.  reads_per_second is
 average value.
 */
function stats_avg_throughput(stats_array)
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
