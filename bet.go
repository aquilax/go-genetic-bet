package main

import(
  "os"
  "strconv"
  "fmt"
  "rand"
  "time"
  "math"
  "sort"
)

const (
  MINBET = 2.
  POPULATION = 250
  RANDOM_SELECT = 0.05
  MUTATE = 0.01
)

var (
  money float64 = 10
  coef Coef
)

type Coef [3]float64;

type Individual struct {
  bets Coef
  risk_max float64
  risk_min float64
  fitness float64
}

type Generation []Individual

func (g Generation) Len() int {
  return len(g)
}

func (g Generation) Less(i, j int) bool {
  return math.Fabs(g[i].fitness) < math.Fabs(g[j].fitness)
}

func (g Generation) Swap(i, j int) {
  g[i], g[j] = g[j], g[i]
}

func (g Generation) Print() {
  for _, i := range g {
    fmt.Println(i.fitness);
  }
}

func ReadCoef() (Coef){
  var c Coef
  c[0], _ = strconv.Atof64(os.Args[1]);
  c[1], _ = strconv.Atof64(os.Args[2]);
  c[2], _ = strconv.Atof64(os.Args[3]);
  return c
  //return {(1/c1)*100, (1/c2)*100, (1/c3)*100};
}

func CoefPercent(c Coef) (Coef) {
  for n, _ := range c {
    c[n] = (1/c[n])*100;
  }
  return c
}

func GenIndividual()(Individual) {
  var i Individual
  i.bets[0] = (rand.Float64() * money)
  if i.bets[0] < MINBET {
    i.bets[0] = 0
  }
  i.bets[1] = (rand.Float64() * (money-i.bets[0]))
  if i.bets[1] < MINBET {
    i.bets[1] = 0
  }
  i.bets[2] = rand.Float64() * (money - (i.bets[0]+i.bets[1]))
  if i.bets[2] < MINBET {
    i.bets[2] = 0
  }
  return i
}

func GenPopulation(n int) (Generation) {
  var gen Generation
  for i := 0; i < n; i++ {
    gen = append(gen, GenIndividual())
  }
  return gen
}

func Fitness(i *Individual, coef Coef){
  var outc Coef
  sum := Sum((*i).bets)
  outc[0] = (coef[0] * (*i).bets[0]) - sum
  outc[1] = (coef[1] * (*i).bets[1]) - sum
  outc[2] = (coef[2] * (*i).bets[2]) - sum
  (*i).risk_max = Max(outc);
  (*i).risk_min = Min(outc);

/*  
  (*i).fitness = ((*i).risk_max) + (*i).risk_min
  if (*i).risk_min < -5 {
    (*i).fitness += -100
  }
*/
  (*i).fitness = (*i).risk_min
}

func (i *Individual) Fix (){
  for n, _ := range(i.bets){
    if ((*i).bets[n]) < 2 {
      if (rand.Float64() > .6) {
        (*i).bets[n] = 2
      } else {
        (*i).bets[n] = 0
      }
    }
  }
}

func Grade(g *Generation) float64 {
  var sum = 0.;
  for _, i := range (*g) {
    sum += i.fitness
  }
  return sum / float64(len(*g))
}

func Max(c Coef) float64 {
  var r float64 = math.Inf(-1)
  for n, _ := range c {
    if c[n] > r {
      r = c[n]
    }
  }
  return r
}

func Min(c Coef) float64 {
  var r float64 = math.Inf(1)
  for n, _ := range c {
    if c[n] < r {
      r = c[n]
    }
  }
  return r
}

func Sum(c Coef) float64 {
  return c[0] + c[1] + c[2];
}

func (i Individual) Mutate() {
  for n ,_ := range i.bets {
    i.bets[n] += rand.Float64()-.5;
  }
}

func Mate(m, f Individual) (Individual) {
  var r Individual
  r.bets[0] = (f.bets[0] + m.bets[0])/2
  r.bets[1] = (f.bets[1] + m.bets[1])/2
  r.bets[2] = (f.bets[2] + m.bets[2])/2
  r.Fix()
  return r
}

func (g Generation) Evolve(retain int) (Generation){
  var ng Generation
  //Add some random loosers for diversity
  for i := retain; i < g.Len(); i++ {
    if RANDOM_SELECT > rand.Float64() {
      g[i].Fix()
      ng = append(ng, g[i])
    }
  }
  for i := 0; i < retain; i++ {
    if MUTATE > rand.Float64() {
      g[i].Mutate()
      g[i].Fix()
    }
  }
  desired_len := len(g) - len(ng)
  for i := 0; i < desired_len; i++ {
    //Mate
    i1 := int(rand.Float64() * float64(retain));
    i2 := int(rand.Float64() * float64(retain));
    ng = append(ng, Mate(g[i1], g[i2]))
  }
  return ng
}


func GenerationNext(gen Generation, generations int, coef Coef) {
  for g := 1; g <= generations; g++ {
    for i, _ := range gen {
      Fitness(&gen[i], coef);
    }
    sort.Sort(gen);
    fmt.Printf("\n== Generation %d == TOP5 == Grade %0.2f =================\n", g, Grade(&gen))
    for i := 0; i  <5; i++ {
      gen[i].Print()
    }
    gen = gen.Evolve(POPULATION/2)
  }
}

func (i Individual) Print(){
  fmt.Printf("risk: %0.2f\t%0.2f\t%0.2f\t", i.risk_max, i.risk_min, i.fitness);
  fmt.Println(i.bets);
}

func main () {
  rand.Seed(time.Nanoseconds());
  coef = ReadCoef();
  perc := CoefPercent(coef)
  fmt.Printf("%0.2f\t%0.2f\t%0.2f\t=%0.2f\n", perc[0], perc[1], perc[2], Sum(perc))

  init := GenPopulation(POPULATION)

  GenerationNext(init, 15, coef);
/*
  var best Individual;
  var bestf float64 = 1000
  for i, _ := range population {
    Fitness(&population[i], coef);
    if math.Fabs(population[i].fitness) < bestf && math.Fabs(population[i].fitness) > 0 {
      best = population[i]
      bestf = math.Fabs(population[i].fitness)
    }
  }
  sort.Sort(population);
  population.Print();
  population = population.Evolve(10)
  fmt.Println(best,"\n")
*/
}
