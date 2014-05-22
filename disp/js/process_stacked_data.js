
/**
 @param {float} throughput
 @param {int} num_threads
 @param {StackedSubData}
 */
function StackedRun(throughput, num_threads,averaged_stacked_sub_data)
{
    this.throughput = throughput;
    this.num_threads = num_threads;
    this.averaged_stacked_sub_data = averaged_stacked_sub_data;
}

/**
 @param {int} time --- The time it took for this operation to complete

 @param {string} label --- The name of this operation

 Note: with children, has a tree structure.
 */
function StackedSubData (time,label)
{
    this.time = time;
    this.label = label;
    // Each element is a StackedSubData object.
    this.children = [];

    // get set in update_cumulatives.  corresponds to when in
    // execution this started and ended.
    this.start = null;
    this.end = null;
}

/**
 @param {StackedSubData} child_to_add
 */
StackedSubData.prototype.add_child = function(child_to_add)
{
    this.children.push(child_to_add);
};

/**
 Returns the maximum depth of this subdata tree.
 */
StackedSubData.prototype.max_depth = function()
{
    var max_child_depth = 0;
    for (var index in this.children)
    {
        var child_depth = this.children[index].max_depth();
        if (child_depth > max_child_depth)
            max_child_depth = child_depth;
    }

    return 1 + max_child_depth;
};

function FlattenedData(label, depth, start, end)
{
    this.label = label;
    this.depth = depth;
    this.start = start;
    this.end = end;
}
FlattenedData.prototype.debug_print = function ()
{
    console.log(
        this.label + " " + this.depth + "; " + this.start +
            "-" + this.end);
};



/**
 @param {int} depth_offset --- How many levels down tree we are.
 
 @param {int} flatten_to_depth --- Total depth of tree.  If current
 node terminates, then return
 
   (flatten_to_depth - depth_offset)

 copies of this node with different depths.
 
 @returns {list} --- Each element is a FlattenedData object.

 Example:
      a
     / \
    b  c
      / \
     d   e

 [
    a: depth 1,
    b: depth 2,
    b: depth 3,
    c: depth 2,
    d: depth 3,
    e: depth 3
 ]
 */
StackedSubData.prototype.flatten = function(depth_offset,flatten_to_depth)
{
    var to_return = [];
    var max_depth = this.max_depth();

    to_return.push(
        new FlattenedData(this.label,depth_offset,this.start,this.end));
    
    if (this.children.length == 0)
    {
        // this is a leaf node.  we need to fill in to_return with
        // copies of self
        for (var i = depth_offset+1; i < flatten_to_depth; ++i)
        {
            to_return.push(
                new FlattenedData(this.label,i,this.start,this.end));
        }
    }
    
    for (var index in this.children)
    {
        var child = this.children[index];
        to_return =
            to_return.concat(child.flatten(flatten_to_depth,depth_offset+1));
    }
    return to_return;
};

/**
 Instead of keeping track of the absolute time deltas, also keep track
 of when, relative to other steps, this sub data ran.
 */
StackedSubData.prototype.update_cumulatives = function(start_offset)
{
    this.start = start_offset;
    this.end = start_offset + this.time;

    var prev_child_end_time = this.start;
    for (var index in this.children)
    {
        var child = this.children[index];
        child.update_cumulatives(prev_child_end_time);
        prev_child_end_time = child.end;
    }
};

/**
 @param {list} sub_data_list --- Each element is a StackedSubData
 object.  Each StackedSubData object should have the same structure.
 Run through all and average runtimes for each.
 */
function average_stacked_sub_data_list(sub_data_list)
{
    // First find average runtime of top nodes.  Then, go through all
    // nodes' children and find average runtime of those, etc.
    var total_node_runtime = 0;
    var label = null;
    var num_children = null;
    for (var index in sub_data_list)
    {
        var sub_data = sub_data_list[index];
        if (label === null)
            label = sub_data.label;
        if (num_children === null)
            num_children = sub_data.children.length;
        
        // sanity check to ensure averaging across same elements
        if (label != sub_data.label)
            throw 'Mismatched labels when averaging';
        if (num_children != sub_data.children.length)
            throw 'Mismatched number of children';
        // end sanity check
        
        total_node_runtime += sub_data.time;
    }
    var avg_node_runtime = total_node_runtime / sub_data_list.length;
    var averaged_node = new StackedSubData(avg_node_runtime,label);

    // go through children and create an average child for the average
    // node to return.
    var averaged_children = [];
    for (var child_index=0; child_index < num_children; ++child_index)
    {
        // each element is a StackedSubData object corresponding to
        // the child of a sub_data_list element with index
        // child_index.
        var child_index_children = [];
        for (var sub_data_index in sub_data_list)
        {
            var sub_data = sub_data_list[sub_data_index];
            child_index_children.push(sub_data.children[child_index]);
        }

        var averaged_child = average_stacked_sub_data_list(child_index_children);
        averaged_node.add_child(averaged_child);
    }
    return averaged_node;
}




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
        var averaged_stacked_sub_data =
            average_stacked_sub_data_list(list_stacked_sub_data);
        averaged_stacked_sub_data.update_cumulatives(0);
        
        var stacked_run = new StackedRun(
            read_only_result.ops_per_second,
            read_only_result.num_threads,
            averaged_stacked_sub_data);
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