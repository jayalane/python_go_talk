
                        if True:  # not _PROXY_SEND_BAD_DATA:
                            if '"X-PAYPAL-OPERATION-NAME=-" "X-PAYPAL-API-RC=-" "ORIG-URL=/cgi-bin/webscr"' in m:
                                ml.ld('skipping webscr message {}', m)
                                ts.incr_data_for_sherlock("socket_msg_filtered", 1, _dim_name)
                                continue
                            if '"ORIG-URL=/akamai/sureroute-test-object.html"' in m:
                                ml.ld('skipping Akamai message {}', m)
                                ts.incr_data_for_sherlock("socket_msg_filtered", 1, _dim_name)
                                continue
                            if 'apache_access' in m and 'Paypal-Debug-Id' not in m:
                                ml.ld('skipping message without debug id but with apache {}', m)
                                ts.incr_data_for_sherlock("socket_msg_filtered", 1, _dim_name)
                                continue
                            t = m.split(':')[0]  # these are host name based
                            if 'apache_access' in t and 'slingshotsiloapi' in t:
                                ml.ld('skipping message from silo for apache {}', t)
                                ts.incr_data_for_sherlock("socket_msg_filtered", 1, _dim_name)
                                continue
                            if 'apache_access' in t and 'shotapi' not in t:
                                ml.ld('skipping message from web for apache {}', t)
                                ts.incr_data_for_sherlock("socket_msg_filtered", 1, _dim_name)
                                continue
                            if 'apache_access' in t and 'slcsb' in t:
                                ml.ld('skipping message from SB for apache {}', t)
                                ts.incr_data_for_sherlock("socket_msg_filtered", 1, _dim_name)
                                continue
                            if ('"url_path":"/tealeaftarget"' in m) or ('"url_path":"/tealeaftarget/"' in m):
                                ml.ld('skipping message from tealeaftarget')
                                ts.incr_data_for_sherlock("socket_msg_filtered", 1, _dim_name)
                                continue
