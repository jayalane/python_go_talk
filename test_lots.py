_TEST_GLETS = []

# Method for creating multiple greenlets for testing
def worker_test(num_greenlets):
    global _TEST_GLETS
    ml.la("{}", num_greenlets)
    test_leader = Leadership('test_host_' + str(num_greenlets), 'ALL_COLOS')
    _TEST_GLETS[num_greenlets] = test_leader.check_leader
    while 1:
        gevent.sleep(2)
        ml.la("{}: Am I Leader: {}", num_greenlets, test_leader.am_i_leader())

# Main function used to test the Leadership algorithm.
if __name__ == "__main__":
    _INACTIVITY_TIME_LIMIT = 2
    _TIME_INTERVAL_TO_CHECK = 2
    glets = []
    for x in xrange(10):
        ml.la("Starting {}", x)
        glets.append(gevent.spawn(worker_test, x))
        _TEST_GLETS.append(None)
    while 1:
        print "Main Loop"
        n = random.randint(0, len(glets) - 1)
        ml.la("Killing greenlet {}", n)
        glets[n].kill()
        try:
            _TEST_GLETS[n].kill()
        except:
            pass
        if random.random() < 0.9:
            gevent.sleep(10)
        glets[n] = gevent.spawn(worker_test, n)
        gevent.sleep(.5)
