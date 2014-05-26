STACKED_EVENT_MAP_DIV_ID = 'stacked_event_map_tests';
STACKED_EVENT_MAP_NOTES_DIV_ID = 'stacked_event_map_tests_notes';

/**
 [{
   "read_only_result":
       {"num_reads": 10000,"num_threads": 1,"op_type": 0,"ops_per_second": 24241.44,"perf_output": null},
   "traces":
       [
           [{"timestamp": 1203040664151479,"event_string": "Creation"},{"timestamp": 1203040664183693,"event_string": "re_ie top"},{"timestamp": 1203040664191778,"event_string": "re_lock top"},{"timestamp": 1203040664197140,"event_string": "re_bottom top"},{"timestamp": 1203040664206224,"event_string": "b_c_root top"},{"timestamp"
: 1203040664211186,"event_string": "c_iter top"},{"timestamp": 1203040664216691,"event_string": "c_iter bottom"},{"timestamp": 1203040664220761,"event_string": "r_and_p top"},{"timestamp": 1203040664296720,"event_string": "r_and_p bottom"},{"timestamp": 1203040664306890,"event_string": "b_c_root bottom"},{"timestamp": 1203040664312134,"event_string": "re_ie bottom"},{"timestamp": 1203040664317860,"event_string": "end_sentinel"}],

         ...
       ]},
    ...
 ]
 */
function event_map_stacked(list_event_map_data)
{
    // Keys are ints (num threads).  Values are dicts, whose keys are
    // strings (first event name) and values are lists of traces with
    // this first event. Eg.,
    //
    // {
    //     1: { event_name: [trace1, trace2, ...]},
    //     2: { event_name: [trace1, trace2, ...]}
    // }
    var unique_trace_names_dict = {};
    for (var i=0; i < list_event_map_data.length; ++i)
    {
        var event_map_data = list_event_map_data[i];
        var num_threads = event_map_data.read_only_result.num_threads;
        // keys are strings (second event name) and values are lists of
        // traces with this second event.
        var num_thread_unique_trace_names_dict = {};
        unique_trace_names_dict[num_threads] =
            num_thread_unique_trace_names_dict;

        for (var trace_index=0; trace_index < event_map_data.traces.length;
             ++trace_index)
        {
            var trace = event_map_data.traces[trace_index];
            // remove final trace from all (had added an un-needed
            // 'end_sentinel' token at end of logging that gets in way of
            // reusing code in process_stacked_data.
            trace = trace.splice(0,trace.length -1);
            // all events begin with name, "Creation," but may have
            // different second values.
            var second_event_name = trace[1].event_string;
            if (!(second_event_name in num_thread_unique_trace_names_dict))
                num_thread_unique_trace_names_dict[second_event_name] = [];
            num_thread_unique_trace_names_dict[second_event_name].push(trace);
        }
    }

    // produce subdata lists, average them, and turn them into a list
    // of StackedRun objects.
    var stacked_run_list = [];
    for (var num_threads in unique_trace_names_dict)
    {
        var event_name_statistics = unique_trace_names_dict[num_threads];
        for (var event_name in event_name_statistics)
        {
            var trace_list = event_name_statistics[event_name];
            var stacked_sub_data_list = process_traces(trace_list);
            var averaged_stacked_sub_data =
                average_stacked_sub_data_list(stacked_sub_data_list);
        
            var stacked_run = new StackedRun(
                0,num_threads, averaged_stacked_sub_data);
            var total_time_ns = stacked_run.averaged_stacked_sub_data.time;
            var total_time_us = total_time_ns / 1000;
            
            stacked_run.averaged_stacked_sub_data.label =
                'Total ' + total_time_us.toFixed(2) + 'us.  ' + num_threads + ' threads.';
            stacked_run_list.push(stacked_run);
        }
    }

    // display stacked run lists
    draw_stacked_graphs(
        stacked_run_list,STACKED_EVENT_MAP_DIV_ID,
        STACKED_EVENT_MAP_NOTES_DIV_ID);
}

