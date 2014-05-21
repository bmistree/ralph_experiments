
/**
 @param {float} throughput
 @param {int} num_threads
 @param {StackedSubData}
 */
function StackedRun(throughput, num_threads,list_stacked_sub_data)
{
    this.throughput = throughput;
    this.num_threads = num_threads;
    this.list_stacked_sub_data = list_stacked_sub_data;
}


/**
 @param {int} time --- The time it took for this operation to complete

 @param {string} label --- The name of this operation
 */
function StackedSubData (time,label)
{
    this.time = time;
    this.label = label;
    // Each element is a StackedSubData object.
    this.children = [];
}

/**
 @param {StackedSubData} child_to_add
 */
StackedSubData.prototype.add_child = function(child_to_add)
{
    this.children.push(child_to_add);
};


/**
 @param {list} data_list --- Each element is an object of the form:
 {
     read_only_result: {
         num_reads: <int>,
         num_threads: <int>,
         ...
         traces: [ @see process_traces for format of traces]
     }
 }

 @returns{list} --- Each element is a StackedRun object.
 */
EXAMPLE = null;
function process_stacked_data(data_list)
{
    var to_return = [];
    for (var index in data_list)
    {
        var datum = data_list[index];
        var read_only_result = datum.read_only_result;
        var list_stacked_sub_data = process_traces(datum.traces);
        var stacked_run = new StackedRun(
            read_only_result.ops_per_second,
            read_only_result.num_threads,list_stacked_sub_data);
        to_return.push(stacked_run);
    }
    EXAMPLE = to_return;
    return to_return;
}

/**
 @param {list} traces_list --- Each element is an array of the following
 form:
 
 [ {timestamp: number, event_string: text},
 
   {timestamp: number, event_string: text 1 top},
   {timestamp: number, event_string: text 1 bottom},
 
   {timestamp: number, event_string: text 2 top},
   {timestamp: number, event_string: text 3 top},
   {timestamp: number, event_string: text 4 top},
   {timestamp: number, event_string: text 4 bottom},
   {timestamp: number, event_string: text 3 bottom},
   {timestamp: number, event_string: text 2 bottom}]

 Ie., the first element text does not have the word top/bottom in it
 and has no match at the end.  All other entries have a matching top
 or bottom.  bottom comes after top, but can have arbitrarily many
 matching entries between its top and bottom.  

 @returns {list} --- Each element is a StackedSubData object
 */
function process_traces(trace_list)
{
    var to_return = [];
    for (var index in trace_list)
    {
        var single_trace = trace_list[index];
        to_return.push(process_single_trace(single_trace));
    }
    return to_return;
}

/**
 @param {trace_list} --- @see process_traces
 
 @returns {StackedSubData} --- label should be text.  time should be
 times between timestamps.  Children should be all top-bottoms
 between.
 */
function process_single_trace(trace_list)
{
    // generate the time it took to generate the first entry
    var time_first = trace_list[0].timestamp;
    var time_last = trace_list[trace_list.length -1].timestamp;
    var total_time = time_last - time_first;
    // using unique label for first event, Total.
    var to_return = new StackedSubData(total_time,"Total");

    var children_array = create_sub_data_children(trace_list.slice(1));
    for (var index in children_array)
        to_return.add_child(children_array[index]);
    return to_return;
}


/**
 @param {list} trace_list --- @see process_stacked_data, except without
 first entry.

 @returns {list} --- Each element 
 */
function create_sub_data_children(trace_list)
{
    var to_return  = [];

    // Basic algorithm, run through list until we find entry that ends
    // in "top".  Then, collect all entries between top and bottom and
    // call create_sub_data_children on them.
    var found_top_looking_for_bottom = false;
    // the entry whose bottom has to match
    var top_prefix = '';
    var top_timestamp = -1;
    // entries between top and bottom
    var intermediate_entries = [];
    for (var index in trace_list)
    {
        var datum = trace_list[index];

        if (found_top_looking_for_bottom)
        {
            if (bottom_with_prefix(datum.event_string,top_prefix))
            {
                // generate subdata item
                var bottom_timestamp = datum.timestamp;
                var total_time = bottom_timestamp - top_timestamp;
                var sub_data = new StackedSubData(total_time,top_prefix);
                var sub_data_children =
                    create_sub_data_children(intermediate_entries);
                for (var child_index in sub_data_children)
                    sub_data.add_child(sub_data_children[child_index]);
                to_return.push(sub_data);

                // reset state machine
                found_top_looking_for_bottom = false;
                intermediate_entries = [];
            }
            else
            {
                // must be an intermediate entry
                intermediate_entries.push(datum);
            }
        }
        else
        {
            // we have not yet found a top entry.  We know that if we
            // have not found a top entry, the next entry *must* be a
            // top entry.  However, for sanity, actually checking this.
            if (! is_top_entry(datum.event_string))
                throw 'No top entry found';
            
            top_timestamp = datum.timestamp;
            top_prefix = extract_top_prefix(datum.event_string);
            found_top_looking_for_bottom = true;
        }
        
    }
    // if we found a top with an unmatched bottom, something's screwy
    // with the data: throw an error
    if (found_top_looking_for_bottom)
        throw 'No bottom entry found for top with prefix ' + top_prefix;

    return to_return;
}


TOP_SUFFIX = 'top';
TOP_SUFFIX_LENGTH = TOP_SUFFIX.length;
BOTTOM_SUFFIX = 'bottom';


/**
 @param {String} event_string --- Should be of the form:

 obj complete_commit top
 Add touched top
 etc.

 Ie., a string that ends with the word "top."

 @returns{string} --- The string without the ending "top".
 */
function extract_top_prefix(event_string)
{
    // debug
    if (! is_top_entry(event_string))
    {
        throw 'Cannot extract top prefix for event ' +
            'that does not end in "top."';
    }
    // end debug
    
    return event_string.substr(0,event_string.length - TOP_SUFFIX_LENGTH);
}

/**
 @returns {boolean} --- true if event_string contains top suffix, false
 otherwise.
 */
function is_top_entry(event_string)
{
    return event_string.indexOf(TOP_SUFFIX) != -1;
}


/**
 @returns {boolean} --- true if event_string is a bottom entry and
 contains top_prefix.
 */
function bottom_with_prefix(event_string,top_prefix)
{
    // check is bottom first
    if (event_string.indexOf(BOTTOM_SUFFIX) == -1)
        return false;

    if (event_string.indexOf(top_prefix) == -1)
        return false;

    return true;
}