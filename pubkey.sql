BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS "PubKeys" (
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  Email TEXT NOT NULL,
  PubKey TEXT
);
CREATE UNIQUE INDEX IF NOT EXISTS PubKeys_ID_IDX ON PubKeys(ID, Email);
INSERT INTO PubKeys(Email, PubKey) VALUES('111174663967311@PHONEWORD.com',replace('-----BEGIN PGP PUBLIC KEY BLOCK-----\n\nmQENBF/ITVUBCACmfD1LXgZqNqyk9T0nvkLP/wjTBW5fhmKYH9ndrJvCj401CCnt\nCdyxeM7orJCdUed/OyXZsb2FMhkJlAXO2KGcH3uRiJRwIEwa+u02eFkUd9HfTsJo\nIj8T+NNv3VrQQCwPUl0g8yKytq3Bn3NUrEHCULACmP2z6qdNJs8RTvUMObFunu0O\ndnhoSbKac5N4x5gAGcyER6O9oWjw4tSicCn/2AADXHlrMc7q3cD4NfwS4CQm84zT\n/nvQ4H4WWP8mj66piubttclVLQuQiUG44pzAW/sZV943owCIsTpkAwRvakce6QMt\np0FmHl/MBblzTcBfrIOuECp8tjbHEI5N9NCVABEBAAG0YFBIT05FV09SRDExMTE3\nNDY2Mzk2NzMxMSAoVGhpcyBpcyBhIHRlc3QgcHVibGljIGFuZCBwcml2YXRlIGtl\neSkgPDExMTE3NDY2Mzk2NzMxMUBQSE9ORVdPUkQuY29tPokBTgQTAQoAOBYhBJ8A\nAMr033NGaOgySUwVHm/cieVZBQJfyE1VAhsDBQsJCAcCBhUKCQgLAgQWAgMBAh4B\nAheAAAoJEEwVHm/cieVZiagH/1qT3Mdt97jEDYAlPw4SI6yVHE7g7lZ/ByxBfH43\nu2z9O0M9ky1iIwMMD57sB2RzccjC8r5HJHP4SW9TE7WwUld7HPWWuq2S/29yBZSv\nV2B06Ubb1I8EjC6KBakKu0yXdnbpZ7WLt0/RLM6k3E+hCqi+tzrA5OkxBnoRgpHn\n468OB/BGmj8YZKMnQP3u9Aist0iLDfuqn5vAfRtoM5zpeyIOc1AvvDuURwh0cNtS\n5t8RV9WpeJFXu/p19nn1DUAwkJcecLHvL9h4uxHZTo4FiMA9JKlAr0u4CT0t8fPI\ns/5fDLbhEE9MZxs/2pSf/rMHdNi2Rla3VsJkk+AhyoBCYCG5AQ0EX8hNVQEIAK6B\nhjxFvw/WAqJnJWBxf0Z21xJFq5h+k15Jqzz2BWCuZZwyocPhs2IXu9e6pcbeakoV\nXPRN0PwHAXC1i+Dzu0qnuIV34XgEggoclKhkF4nkjRx9AZj53oG3fwbzVx2X3pBF\n98bZlmXBSiGCikA1TZgxHBPfvoQJI8YbvOvg5sTq+zEzH0VNkYcgsFbEwAXcgasF\nb8IW2wmh2T8Pso6hOS1jPCTzDfTze24A9b0i1tbfcoTWH7OfrWGoBS+Murrx3xYI\nX8Gr+HmcdvJsCM0ttV2oFOe+efqikUZzLmRcsjBdoe0CdRbncxSzbyEdeQhW6LDs\nkDFsU2tuEOtKL3wf0iEAEQEAAYkBNgQYAQoAIBYhBJ8AAMr033NGaOgySUwVHm/c\nieVZBQJfyE1VAhsMAAoJEEwVHm/cieVZ2bQH/0kRApyivwFJD7V3mipsdoh+21Jk\nAvsmAMZiyY84LEdAGYUltgZddqlTzoo23R/V6mFVVuscdtac4Su2ZGsH9JNUKvvz\nLy1ku3oEIgSDwzqtkUPTq8MLq4C+/Fz87NxYX3znxfkWwKfmgxdzVtSNYW/8AhEH\nnT+cRrZi98vRbk21TiWcBnuFsu96h58U2IQUsxszT5mUIIxYysPhQlptDd2Zw0bD\n89vdYCWCKqwIgIyRRqLw9vVdj2FHWyP/9Ix1+Jq5aqFjncPPmjqo1mdAb0WL722T\nRXnw6oYLZj5efRTkJ/z9yiIUFT4m4o6zUv0bbQoy3y7h65N56xPfk3ZTuvA=\n=l4jS\n-----END PGP PUBLIC KEY BLOCK-----\n\n-----BEGIN PGP PRIVATE KEY BLOCK-----\n\n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\nuvA=\n=2ndq\n-----END PGP PRIVATE KEY BLOCK-----','\n',char(10)));
INSERT INTO PubKeys(Email, PubKey) VALUES('ashuidesigns@gmail.com',replace('-----BEGIN PGP PUBLIC KEY BLOCK-----\n\n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\nwhrwiFjQ5Qv9RNpW/3bYjog8HZ5BCw==\n=r4yu\n-----END PGP PUBLIC KEY BLOCK-----','\n',char(10)));
INSERT INTO PubKeys(Email, PubKey) VALUES('eylanderchad@gmail.com',replace('-----BEGIN PGP PUBLIC KEY BLOCK-----\n\n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\n8R8pYzjtv1L+litoLKAmZoitfM32XUJ4ejalmZo=\n=Xshp\n-----END PGP PUBLIC KEY BLOCK-----','\n',char(10)));
INSERT INTO PubKeys(Email, PubKey) VALUES('jaycarpenter@desertblockchain.com',NULL);
INSERT INTO PubKeys(Email, PubKey) VALUES('jaycarpenter@PHONEWORDworkshops.com',NULL);
INSERT INTO PubKeys(Email, PubKey) VALUES('jose@tafla.com',replace('-----BEGIN PGP PUBLIC KEY BLOCK-----\n\n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\n2njoyQ3R+GaIV2X69GGaPsii\n=KwFf\n-----END PGP PUBLIC KEY BLOCK-----','\n',char(10)));
INSERT INTO PubKeys(Email, PubKey) VALUES('louloucarpentertest1@gmail.com',replace('-----BEGIN PGP PUBLIC KEY BLOCK-----\n\nmI0EX63RcAEEAMkozqYh0WZbdJZX65vd5afnuO0DKbi690LkZ9lhDXkYPx5jKfRU\nJ8ZK8/JWeotDEAvN5Bb6BNIMRKkjek/38gY4KzM9gUt6dlppC4J7YcxAEYjP6eG7\nXfhdgwmGgfwRPIoj6EPR5u7wjCxlaFnNZdlOs+JmCgaYJSQRXwhiNJALABEBAAG0\nKmxvdSB0ZXN0MSA8bG91bG91Y2FycGVudGVydGVzdDFAZ21haWwuY29tPojOBBMB\nCgA4FiEEjjD4o01+qcodsiv7zBmB11eyuXoFAl+t0XACGwMFCwkIBwIGFQoJCAsC\nBBYCAwECHgECF4AACgkQzBmB11eyuXr6NQQAo3x+hXHmWxlsRoD0sqdzjbnBrDTl\n91yE7UA2x7HZNIiCJ7M9fR2BfkpUbS4pSRfW7HJF2gxy4GrEmk3jS/6uU38PYxTc\nQNMEOihg/iTeg9c1GamzGn/LOhSdG7y4PU1bQMdnLGxkwMGyuiGmqXZTEIVj2Onb\nTIO7qFuirsvxBIi4jQRfrdFwAQQAtle2xaB0yuPJUL63lPWTiqYH8bWcPi8oThx5\nx6B7w3IYr2692AP73Op1JQFv83/1Ki+HYaPz4siaJ0Md6nfAuKFriq63O+J0tOS3\nqSbWhaIACLa3sFJNKL6dznz3fzyhbQHiERENOG5G07L9+LeXtbBUIgdALFa1eglK\nb528CIEAEQEAAYi2BBgBCgAgFiEEjjD4o01+qcodsiv7zBmB11eyuXoFAl+t0XAC\nGwwACgkQzBmB11eyuXqorwQAvW4c7oA1sO93Wg3SnGp6Ebs3cKfdO3Rvy/YVdOJg\nmJlInIhHqY5oMIm4lSgKkUEA0F5tXqeD3jg0cvnXgqtfUjRQEnE2BAYDtajnmj55\nSBoo4gqdmMTqvTPcaXKBNf9veNqvCHqnOYImhEnGHctJF0ASeA31iwBTmJzchUna\nwj8=\n=oD3t\n-----END PGP PUBLIC KEY BLOCK-----\n\n\nRSA 1024\n\nroot@raspberrypi:/# ipfs add loulougpgpublickeytest1.txt \nadded QmbFJoSJKNkFzSTsU7Gfp2pvUsndG5yv58VsM9CJ7pa5eU loulougpgpublickeytest1.txt\n 1.04 KiB / 1.04 KiB [=======================================================================================================] 100.00%  \n\n\nPublic Key created 2020-11-12 with no expiration\n\nPHONEWORD Workshops Global Directory Lookup Number:\n\n1111LOULOUCARPENTER\n111156856820000 to 111156856829999\n\nroot@raspberrypi:/# sha256sum loulougpgpublickeytest1.txt \nd3482c8aa080ad906e297620b03049308bb672d70e9c67a69dcfb910390b8153  loulougpgpublickeytest1.txt\n\n-----BEGIN PGP PRIVATE KEY BLOCK-----\n\nlQIGBF+t0XABBADJKM6mIdFmW3SWV+ub3eWn57jtAym4uvdC5GfZYQ15GD8eYyn0\nVCfGSvPyVnqLQxALzeQW+gTSDESpI3pP9/IGOCszPYFLenZaaQuCe2HMQBGIz+nh\nu134XYMJhoH8ETyKI+hD0ebu8IwsZWhZzWXZTrPiZgoGmCUkEV8IYjSQCwARAQAB\n/gcDAuuwrPewGABo5nxOcllLFnXmiIEdhND56VGUq35b+j5MQJ1BDJgEb4A8zrLp\naPSVd1Il2X6v5U+0GRTpKNQyrubjz/1gxZ+sRHCsWAzXQs/+L4VV6jPMWKt3t+ej\nk1PC8JPplIg4yzVGZfD66z09sjouD35OO70zrJVK6Y/aF9o+T77Vel9DtVUG+i1L\nqXz8gp0yY5M3/GMIn2vUC7uTby6HSst8cLT7BDc4ARnd5KbzBkOJB1T8v8IOhGAq\n/Or2FuqTpSS9gOiFsIkrEuZ4Kjutvhy+fh0VqQO7sSJXVVTauDiLQ3IaG5IkhLsH\nLw79zH2X8+mMtk1tmPyMmPMpk4GPvzLQ6jJTlyXopbOhtkB+CStGcdzYts3Kp8OZ\nQr9rz8HWwrZ66HlIMvgLlj2BeTTfPOa2j3oRorGcbX1cSIDzwR7JEgcnqlIUnXHa\nR+Sp2V1tlX5zXcgFE8BmaQivLOPoWp8oHRjFzNbnEuyRb4duKy/d5Fy0KmxvdSB0\nZXN0MSA8bG91bG91Y2FycGVudGVydGVzdDFAZ21haWwuY29tPojOBBMBCgA4FiEE\njjD4o01+qcodsiv7zBmB11eyuXoFAl+t0XACGwMFCwkIBwIGFQoJCAsCBBYCAwEC\nHgECF4AACgkQzBmB11eyuXr6NQQAo3x+hXHmWxlsRoD0sqdzjbnBrDTl91yE7UA2\nx7HZNIiCJ7M9fR2BfkpUbS4pSRfW7HJF2gxy4GrEmk3jS/6uU38PYxTcQNMEOihg\n/iTeg9c1GamzGn/LOhSdG7y4PU1bQMdnLGxkwMGyuiGmqXZTEIVj2OnbTIO7qFui\nrsvxBIidAgYEX63RcAEEALZXtsWgdMrjyVC+t5T1k4qmB/G1nD4vKE4cecege8Ny\nGK9uvdgD+9zqdSUBb/N/9Sovh2Gj8+LImidDHep3wLiha4qutzvidLTkt6km1oWi\nAAi2t7BSTSi+nc589388oW0B4hERDThuRtOy/fi3l7WwVCIHQCxWtXoJSm+dvAiB\nABEBAAH+BwMCPzfIB9oaS57mvN6xYcR7YaU2fFM4xJL9FYCLbTetGwwtItdVwN1P\nXRlLjhy0LX62GrWMA6nI28m/Z67/s+AGHXbBfPy7PIT3HwnR9qFplvqy946DuweH\nDxt/ssmww2rSK8PxX2CJbk7UqLDcVMqX4iWkPX6nUnTYM8yloGAZObxg837/xxJU\nRxNpoYR7dcCObpJYOMK5YCZmxdyaEhe5RmlEgWFawC1Pj/SJN/OLhrT0k4w9eV66\ngd2a5G3vo+7Ed1IFFcVspVTwgOL8iXWdMGT9Z5I0uNFDvmplyY0cSecIfKzGsHj5\nBEGBU1eMRYWuAt392xaoIIFLe0pWNEbWPSKR9iKanUKXxznZe0SedK+ss5/GnzbK\nzS7aiGIf7BjVZB3frnnPfJOaUPdMhwiO1gRL1OrhuvnHhk/igkmLY86A1ekXYQBe\nHXb9wZLefyJsAkozFcy2yHasra443kgAnUJ0YJH/dHXIw0znRxWfCzO+hmubbIi2\nBBgBCgAgFiEEjjD4o01+qcodsiv7zBmB11eyuXoFAl+t0XACGwwACgkQzBmB11ey\nuXqorwQAvW4c7oA1sO93Wg3SnGp6Ebs3cKfdO3Rvy/YVdOJgmJlInIhHqY5oMIm4\nlSgKkUEA0F5tXqeD3jg0cvnXgqtfUjRQEnE2BAYDtajnmj55SBoo4gqdmMTqvTPc\naXKBNf9veNqvCHqnOYImhEnGHctJF0ASeA31iwBTmJzchUnawj8=\n=oac4\n-----END PGP PRIVATE KEY BLOCK-----','\n',char(10)));
INSERT INTO PubKeys(Email, PubKey) VALUES('louloucarpentertest2@gmail.com',replace('-----BEGIN PGP PUBLIC KEY BLOCK-----\n\nmI0EX63uNQEEALLeEbaJGdicYsb49InYPYKUEJaymAfeBnP2Ctzz1I/9ufHWol1B\nFiXkH0FMEEsBDc7HePAWdzVs+Y+wvi8qMDYgiH/WXASMNKe9j1ul+yJOZQhGhm57\nQCqc87M5kq17Y6Nwg3GmMWaR8rXJtvZmnR+kygLng8cupp974XwhCm6zABEBAAG0\nKGxvdWxvdTIgPGxvdWxvdWNhcnBlbnRlcnRlc3QyQGdtYWlsLmNvbT6IzgQTAQoA\nOBYhBAOPWR5PMBD08KcH78gYzDOwqvSxBQJfre41AhsDBQsJCAcCBhUKCQgLAgQW\nAgMBAh4BAheAAAoJEMgYzDOwqvSxV5cD/04MpmZz+Yq593XsNISPF7a7/pz76KUb\niEoHNwTJF+P8YXSHbNv78hRFyDVxz5XxAYyGnBjNe3z6JyW8raIkhVfio23Fqe0x\nAwTpbHTk0zHqwjC364aUd3MevP9uWVeI8mRu6XZGy+qtw70exeTsSzATC/YJnOa3\n9BOTCqwL/PY8uI0EX63uNQEEAMuKrM2Kf6OILeuUJ3OggfMdZN9bk76FzWb4tc/L\nbBU6/xPx19V1ykJ2LxjUEkK9SVB7TPcvrMI2+29FqYm0GTS/ZKIBICi99AzgWhHI\nkj1isrlfrqua2md+1Ad6xPSIh1pwfCnsFwTF/6b/b/sUMhBn4WXuCnUKyHhMhEXM\nBLpxABEBAAGItgQYAQoAIBYhBAOPWR5PMBD08KcH78gYzDOwqvSxBQJfre41AhsM\nAAoJEMgYzDOwqvSxED8D/1Qbzm6Kw7X52kD1v012x7sDqrOemfIbjCWBWdlqTFB2\nKnuYnCX9JyLffyCZDZ7E/fgeI0pa94bLymXJodF9ceSd3f8ok32+WlEXjZNLF2PX\nkKZ8CvBf9fAGbaTo5LeK1PFtdaBgPAdUaYbDNeLz8ZVX5CUE0YsEVABTHMUtHYrm\n=SKvM\n-----END PGP PUBLIC KEY BLOCK-----\n\nRSA 1024\n\nroot@raspberrypi:/# ipfs add loulougpgpublickeytest2.txt \nadded QmPBy4gbMjMtopu2HWWxwV52tX3pYu3tKjACEonot3rUNm loulougpgpublickeytest2.txt\n 1.03 KiB / 1.03 KiB [=======================================================================================================] 100.00%\n\nPublic Key created 2020-11-12 with no expiration\n\nPHONEWORD Workshops Global Directory Lookup Number:\n\n1111LOULOUCARPENTER2\n111156856820000 to 111156856829999\n\nroot@raspberrypi:/# sha256sum loulougpgpublickeytest2.txt \na79fff2cbe78e0512e1fb56349f8d21527224eae27c1c68a9b9f68e1f91c94e1  loulougpgpublickeytest2.txt\n\n-----BEGIN PGP PRIVATE KEY BLOCK-----\n\n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\n=Fd+M\n-----END PGP PRIVATE KEY BLOCK-----','\n',char(10)));
INSERT INTO PubKeys(Email, PubKey) VALUES('louloucarpentertest3@gmail.com',replace('-----BEGIN PGP PUBLIC KEY BLOCK-----\n\nmI0EX63rkAEEAKdhK5qXa8fZkYez5JiCiswIuPr2bgL95LiBl/FXKiKvcylWXVdz\n0/A1wGpR1G3CkuMXig1M5wf59t73XrNus4AKThB4Jh3aEELZALNrbKv2StefPFNf\nyd6Bt1QhF7DBHsUkIwK2x0QI8uDRlfYU/IfiYU1k6oJ6GgvNQ0Wo15fVABEBAAG0\nKGxvdWxvdTMgPGxvdWxvdWNhcnBlbnRlcnRlc3QzQGdtYWlsLmNvbT6IzgQTAQoA\nOBYhBDuA+cdEg4afDwG0UpTdOZYDhCqfBQJfreuQAhsDBQsJCAcCBhUKCQgLAgQW\nAgMBAh4BAheAAAoJEJTdOZYDhCqfnzID/i3r9ERXKb8EPIyKYUooovNpSwukJHoM\n5ig2A7Ai87AgdYJMooNQZ8sVFGljeFOWVi8X8TpoeWKFQCxccNrYx0m+latWXZE7\n1A5oHtD3wnyELYXnb8V+rxrPo/W9eAqYlBRLFWg9Eehbe4GdE2dmkPCkw9sT3TNf\nXpZDzY4dgC8+uI0EX63rkAEEAMdMgyPuY7ZyAhaDycZOEsxbSdHFEqLsnwyhCpXv\n1BXQmHtZmQMoCP1pYtsgU12dEYnIibyibqrH95fFXN6GKMT0Lkk2CnpSonKH8K4P\n1bNOQdQbEqMnhq1a4s7+GaqSWqKqViY/lE325t0gODRugDCCaMS4NFMc9XL5LP3b\nw48ZABEBAAGItgQYAQoAIBYhBDuA+cdEg4afDwG0UpTdOZYDhCqfBQJfreuQAhsM\nAAoJEJTdOZYDhCqfY/gD/RPHU1D85tuBfCjntATnf2WgnxzfJBhlF8EGBsyEQi1J\nRBFHXcAgFfZ6AdrDm0lIzFatMZFg+MrYnd37ZqUa7AayOygK8Mfo2FTItCuykMaf\naO0z5hg20AXi98bRchk0RjrohZfG1ZOh9hYTJVOElGTFtWoc4Z1ZsYQ1OhWzj6lS\n=AQGT\n-----END PGP PUBLIC KEY BLOCK-----\n\n\n-----BEGIN PGP PRIVATE KEY BLOCK-----\n\n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\n=os8r\n-----END PGP PRIVATE KEY BLOCK-----','\n',char(10)));
INSERT INTO PubKeys(Email, PubKey) VALUES('randy.rgc@outlook.com',replace('-----BEGIN PGP PUBLIC KEY BLOCK-----\n\n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\n/sB3UQ8s3M/8kmJqnJAD0Eu5zoghPsbFONhufQ==\n=St1t\n-----END PGP PUBLIC KEY BLOCK-----','\n',char(10)));
INSERT INTO PubKeys(Email, PubKey) VALUES('rjaycarpenter@gmail.com',replace('-----BEGIN PGP PUBLIC KEY BLOCK-----\n\n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\nvW0n0DS8/50iAEGi2C5L58T6lKGGlFfRxKMpRw==\n=kP/J\n-----END PGP PUBLIC KEY BLOCK-----','\n',char(10)));
COMMIT;