GC_OFF_TESTS_DIV_ID = 'gc_off_tests';

/**
 @param {list} data_list --- Each element is a list containing
 integers and floats.  Indices are formatted according to 
 */
function gc_off_graph(data_list)
{
    draw_num_threads_graph(data_list,GC_OFF_TESTS_DIV_ID);
}
