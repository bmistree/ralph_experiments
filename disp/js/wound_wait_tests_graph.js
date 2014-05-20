WOUND_WAIT_TESTS_DIV_ID = 'wound_wait_tests';

/**
 @param {list} data_list --- Each element is a list containing
 integers and floats.  Indices are formatted according to 
 */
function wound_wait_graph(data_list)
{
    draw_num_threads_graph(data_list,WOUND_WAIT_TESTS_DIV_ID);
}
