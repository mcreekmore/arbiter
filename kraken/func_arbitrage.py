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
            if pair_b != pair_a:
                if b_base in a_pair_box or b_quote in a_pair_box:

                    # get pair c
                    for pair_c in pairs:
                        c_base = pair_c['base']
                        c_quote = pair_c['quote']
                        c_pair_box = [c_base, c_quote]

                        # count number of matching c items
                        if pair_c != pair_a and pair_c != pair_b:
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
                                        "unique_item": unique_item
                                    }

                                    triangular_pairs_list.append(match_dict)
    print(triangular_pairs_list[:10])
    print(len(triangular_pairs_list))
