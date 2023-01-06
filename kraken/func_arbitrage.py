import requests


def get_tradeable_pairs(pairs):
    tradeable_pairs = []

    for pair in pairs:
        # add pair name to object
        pairs[pair]["pair_name"] = pair

        # only allows 'online' pairs
        if pairs[pair]['status'] == 'online':
            tradeable_pairs.append(pairs[pair])

    return tradeable_pairs


def structure_arbitrage_pairs(pairs):
    triangular_pairs_list = []
    remove_duplicates_list = []
    pairs_list = pairs[:10]

    # get pair a
    for pair_a in pairs:
        a_base = pair_a['base']
        a_quote = pair_a['quote']
        # assign a to a box
        a_pair_box = [a_base, a_quote]

        # get pair b
        for pair_b in pairs:
            b_base = pair_b['base']
            b_quote = pair_b['quote']
            b_pair_box = [b_base, b_quote]

            # check pair b
            if b_pair_box != a_pair_box:
                if b_base in a_pair_box or b_quote in a_pair_box:

                    # get pair c
                    for pair_c in pairs:
                        c_base = pair_c['base']
                        c_quote = pair_c['quote']
                        c_pair_box = [c_base, c_quote]

                        # count number of matching c items
                        if c_pair_box != a_pair_box and c_pair_box != b_pair_box:
                            combine_all = [pair_a, pair_b, pair_c]
                            pair_box = [a_base, a_quote, b_base,
                                        b_quote, c_base, c_quote]

                            counts_c_base = 0
                            counts_c_quote = 0
                            for i in pair_box:
                                if i == c_base:
                                    counts_c_base += 1
                                if i == c_quote:
                                    counts_c_quote += 1

                            # determine triangular match
                            if counts_c_base == 2 and counts_c_quote == 2 and c_base != c_quote:
                                # combined = sorted(a_pair_box), sorted(
                                #     b_pair_box), sorted(c_pair_box)

                                combined = [a_base + '/' + a_quote, b_base +
                                            "/" + b_quote, c_base + "/" + c_quote]
                                unique_item = sorted(combined)
                                # print(unique_item)

                                # if unique_item == ['AAVE/ZEUR', 'AAVE/ZUSD', 'ZEUR/ZUSD']:
                                #     print('match!')

                                if unique_item not in remove_duplicates_list:
                                    match_dict = {
                                        "a_base": a_base,
                                        "b_base": b_base,
                                        "c_base": c_base,
                                        "a_quote": a_quote,
                                        "b_quote": b_quote,
                                        "c_quote": c_quote,
                                        "a_pair": a_pair_box,
                                        "b_pair": b_pair_box,
                                        "c_pair": c_pair_box,
                                        "a_altname": pair_a['altname'],
                                        "b_altname": pair_b['altname'],
                                        "c_altname": pair_c['altname'],
                                        "unique_item": unique_item
                                    }

                                    triangular_pairs_list.append(match_dict)
                                    remove_duplicates_list.append(unique_item)
    return triangular_pairs_list


def get_price_for_t_pair(t_pair):
    # get order book for pair
    a_altname = t_pair['a_altname']
    b_altname = t_pair['b_altname']
    c_altname = t_pair['c_altname']

    a_resp = requests.get(
        'https://api.kraken.com/0/public/Depth?pair={}'.format(a_altname)).json()
    b_resp = requests.get(
        'https://api.kraken.com/0/public/Depth?pair={}'.format(b_altname)).json()
    c_resp = requests.get(
        'https://api.kraken.com/0/public/Depth?pair={}'.format(c_altname)).json()

    print(a_resp)
    print(b_resp)
    print(c_resp)

    # extract pairs info
    # pair_a =
