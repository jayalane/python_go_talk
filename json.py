@infra.async.timed
def filter_dict(a_dict, list_of_ok_keys=a_list,
                strip_queries=strip_queries,
                prefix_so_far=''):
    # this is doing a lot of copying
    new_dict = {}
    for k, v in a_dict.iteritems():
        infra.async.maybe_sleep()
        if isinstance(v, dict):
            new_v = filter_dict(v, list_of_ok_keys, strip_queries, prefix_so_far + k + '.')
            if len(new_v) > 0:
                new_dict[k] = new_v
        elif prefix_so_far + k in list_of_ok_keys:
            new_v = v
            if k in strip_queries:
                new_v = v.split('?')[0]
                if 'cmd=' in v:
                    new_v = new_v + "/cmd/" + v.split("cmd=")[1].split('&')[0]
            if k in blanket:
                new_v = "BLANKET"
            new_dict[k] = new_v
    return new_dict
