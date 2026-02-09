# Test setup
This tests require a churchtools instance with following groups:

```yaml
testGroup1:
    type: dynamic group
    members:
        - Anna;Mustermann;Musterstraße 1;11111;Musterstadt;weiblich;2008-09-12
        - Max;Beispiel;Beispielweg 2;22222;Musterdorf;männlich;2008-09-24
        - Lisa;Tester;Testallee 3;33333;Mustertal;weiblich;2008-09-29
        - Tom;Probe;Probestraße 4/1;44444;Musterhausen;männlich;2008-09-18
```

```yaml
testGroup2:
    type: dynamic group
    members:
        - Kerstin;Demonstration;Demoweg 5;55555;Musterberg;weiblich;1982-03-25
        - Markus;Vorschau;Vorschauplatz 6;66666;Musterfeld;männlich;1982-03-04
```

```yaml
testGroup3:
    type: dynamic group
    group updates: disabled
    members: none
```

```yaml
testGroup4:
    type: non dynamic
```

